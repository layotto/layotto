package io.mosn.layotto.v1;

import com.google.common.base.Strings;
import com.google.protobuf.Any;
import com.google.protobuf.ByteString;
import io.grpc.Metadata;
import io.grpc.stub.MetadataUtils;
import io.mosn.layotto.v1.config.RuntimeProperties;
import io.mosn.layotto.v1.exceptions.RuntimeClientException;
import io.mosn.layotto.v1.serializer.ObjectSerializer;
import org.slf4j.Logger;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;
import spec.sdk.runtime.v1.domain.invocation.InvokeResponse;
import spec.sdk.runtime.v1.domain.state.*;

import java.io.Closeable;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;

public class RuntimeClientGrpc extends AbstractRuntimeClient {

    private static final String TIMEOUT_KEY = "timeout";
    protected final RuntimeGrpc.RuntimeBlockingStub blockingStub;
    private final Closeable closeable;

    RuntimeClientGrpc(Logger logger,
                      int timeoutMs,
                      ObjectSerializer stateSerializer,
                      Closeable closeable,
                      RuntimeGrpc.RuntimeBlockingStub blockingStub) {
        super(logger, timeoutMs, stateSerializer);
        this.closeable = closeable;
        this.blockingStub = blockingStub;
    }

    @Override
    public String sayHello(String name, int timeoutMillisecond) {
        try {
            // 1. prepare request
            RuntimeProto.SayHelloRequest req = RuntimeProto.SayHelloRequest.newBuilder()
                    .setServiceName(name)
                    .setName(name)
                    .build();

            // 2. invoke
            RuntimeProto.SayHelloResponse response = blockingStub.withDeadlineAfter(timeoutMillisecond, TimeUnit.MILLISECONDS)
                    .sayHello(req);

            // 3. parse result
            return response.getHello();
        } catch (Exception e) {
            logger.error("sayHello error ", e);
            throw new RuntimeClientException(e);
        }
    }

    @Override
    public InvokeResponse<byte[]> invokeMethod(String appId, String methodName, byte[] data, Map<String, String> header, int timeoutMs) {
        try {
            // 1. prepare request
            final ByteString byteString = ByteString.copyFrom(data);
            Any anyData = Any.newBuilder()
                    .setValue(byteString)
                    .build();
            RuntimeProto.CommonInvokeRequest commonReq = RuntimeProto.CommonInvokeRequest.newBuilder()
                    .setMethod(methodName)
                    .setData(anyData)
                    .build();
            RuntimeProto.InvokeServiceRequest invokeReq = RuntimeProto.InvokeServiceRequest.newBuilder()
                    .setId(appId)
                    .setMessage(commonReq)
                    .build();
            // metadata
            Metadata metadata = new Metadata();
            for (Map.Entry<String, String> entry : header.entrySet()) {
                if (entry.getKey() != null && entry.getValue() != null) {
                    Metadata.Key<String> key = Metadata.Key.of(entry.getKey(), Metadata.ASCII_STRING_MARSHALLER);
                    metadata.put(key, entry.getValue());
                }
            }
            Metadata.Key<String> key = Metadata.Key.of(TIMEOUT_KEY, Metadata.ASCII_STRING_MARSHALLER);
            metadata.put(key, Integer.toString(timeoutMs));

            // 2. invoke
            RuntimeProto.InvokeResponse resp = blockingStub.withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS)
                    .withInterceptors(MetadataUtils.newAttachHeadersInterceptor(metadata))
                    .invokeService(invokeReq);

            // 3. parse result
            InvokeResponse<byte[]> result = new InvokeResponse<>();
            result.setContentType(resp.getContentType());
            byte[] bytes = new byte[]{};
            result.setData(bytes);
            if (resp.getData() == null) {
                return result;
            }
            if (resp.getData().getValue() == null) {
                return result;
            }
            if (resp.getData().getValue().toByteArray() == null) {
                return result;
            }
            result.setData(resp.getData().getValue().toByteArray());
            return result;
        } catch (Exception e) {
            logger.error("invokeMethod error ", e);
            throw new RuntimeClientException(e);
        }
    }

    @Override
    public void publishEvent(String pubsubName, String topicName, byte[] data, String contentType, Map<String, String> metadata) {
        try {
            // 1. prepare data
            if (data == null) {
                data = new byte[]{};
            }
            final ByteString byteString = ByteString.copyFrom(data);
            // Content-type can be overwritten on a per-request basis.
            // It allows CloudEvents to be handled differently, for example.
            if (contentType == null || contentType.isEmpty()) {
                contentType = RuntimeProperties.DEFAULT_PUBSUB_CONTENT_TYPE;
            }

            // 2. prepare request
            RuntimeProto.PublishEventRequest.Builder envelopeBuilder = RuntimeProto.PublishEventRequest.newBuilder()
                    .setTopic(topicName)
                    .setPubsubName(pubsubName)
                    .setData(byteString)
                    .setDataContentType(contentType);
            // metadata
            if (metadata != null) {
                envelopeBuilder.putAllMetadata(metadata);
            }
            RuntimeProto.PublishEventRequest req = envelopeBuilder.build();

            // 3. invoke
            blockingStub.publishEvent(req);
        } catch (Exception e) {
            logger.error("publishEvent error ", e);
            throw new RuntimeClientException(e);
        }
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void saveBulkState(SaveStateRequest request) {
        try {
            final String stateStoreName = request.getStoreName();
            final List<State<?>> states = request.getStates();

            // 1. validate
            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }

            // 2. prepare request
            RuntimeProto.SaveStateRequest.Builder builder = RuntimeProto.SaveStateRequest.newBuilder();
            builder.setStoreName(stateStoreName);
            for (State<?> state : states) {
                RuntimeProto.StateItem stateItem = buildStateRequest(state)
                        .build();
                builder.addStates(stateItem);
            }
            RuntimeProto.SaveStateRequest req = builder.build();

            // 3. invoke
            blockingStub.withDeadlineAfter(getTimeoutMs(), TimeUnit.MILLISECONDS)
                    .saveState(req);
        } catch (Exception e) {
            logger.error("saveBulkState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    private <T> RuntimeProto.StateItem.Builder buildStateRequest(State<T> state) throws IOException {
        RuntimeProto.StateItem.Builder stateBuilder = RuntimeProto.StateItem.newBuilder();
        T value = state.getValue();
        // 1. serialize value
        byte[] bytes;
        if (value == null || value instanceof byte[]) {
            bytes = (byte[]) value;
        } else {
            bytes = stateSerializer.serialize(value);
        }
        if (bytes != null) {
            stateBuilder.setValue(ByteString.copyFrom(bytes));
        }

        // 2. etag
        if (state.getEtag() != null) {
            RuntimeProto.Etag etag = RuntimeProto.Etag.newBuilder()
                    .setValue(state.getEtag())
                    .build();
            stateBuilder.setEtag(etag);
        }
        // 3. metadata
        if (state.getMetadata() != null) {
            stateBuilder.putAllMetadata(state.getMetadata());
        }
        // 4. key
        stateBuilder.setKey(state.getKey());
        // 5. StateOptions
        RuntimeProto.StateOptions.Builder optionBuilder = null;
        if (state.getOptions() != null) {
            StateOptions options = state.getOptions();
            optionBuilder = RuntimeProto.StateOptions.newBuilder();
            if (options.getConcurrency() != null) {
                RuntimeProto.StateOptions.StateConcurrency concurrency = getGrpcStateConcurrency(options);
                optionBuilder.setConcurrency(concurrency);
            }
            if (options.getConsistency() != null) {
                RuntimeProto.StateOptions.StateConsistency consistency = getGrpcStateConsistency(options);
                optionBuilder.setConsistency(consistency);
            }
        }
        if (optionBuilder != null) {
            RuntimeProto.StateOptions stateOptions = optionBuilder.build();
            stateBuilder.setOptions(stateOptions);
        }

        return stateBuilder;
    }

    /**
     * Delete a state.
     *
     * @param request Request to delete a state.
     */
    @Override
    public void deleteState(DeleteStateRequest request) {
        try {
            final String stateStoreName = request.getStateStoreName();
            final String key = request.getKey();
            final StateOptions options = request.getStateOptions();
            final String etag = request.getEtag();
            final Map<String, String> metadata = request.getMetadata();

            // 1. validate
            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }
            if ((key == null) || (key.trim().isEmpty())) {
                throw new IllegalArgumentException("Key cannot be null or empty.");
            }

            // 2. prepare request
            RuntimeProto.StateOptions.Builder optionBuilder = null;
            if (options != null) {
                optionBuilder = RuntimeProto.StateOptions.newBuilder();
                if (options.getConcurrency() != null) {
                    RuntimeProto.StateOptions.StateConcurrency concurrency = getGrpcStateConcurrency(options);
                    optionBuilder.setConcurrency(concurrency);
                }
                if (options.getConsistency() != null) {
                    RuntimeProto.StateOptions.StateConsistency consistency = getGrpcStateConsistency(options);
                    optionBuilder.setConsistency(consistency);
                }
            }
            RuntimeProto.DeleteStateRequest.Builder builder = RuntimeProto.DeleteStateRequest.newBuilder()
                    .setStoreName(stateStoreName)
                    .setKey(key);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            if (etag != null) {
                RuntimeProto.Etag value = RuntimeProto.Etag.newBuilder()
                        .setValue(etag)
                        .build();
                builder.setEtag(value);
            }
            if (optionBuilder != null) {
                RuntimeProto.StateOptions stateOptions = optionBuilder.build();
                builder.setOptions(stateOptions);
            }
            RuntimeProto.DeleteStateRequest req = builder.build();

            // 3. invoke
            blockingStub.withDeadlineAfter(getTimeoutMs(), TimeUnit.MILLISECONDS)
                    .deleteState(req);
        } catch (Exception e) {
            logger.error("deleteState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    private RuntimeProto.StateOptions.StateConcurrency getGrpcStateConcurrency(StateOptions options) {
        switch (options.getConcurrency()) {
            case FIRST_WRITE:
                return RuntimeProto.StateOptions.StateConcurrency.CONCURRENCY_FIRST_WRITE;
            case LAST_WRITE:
                return RuntimeProto.StateOptions.StateConcurrency.CONCURRENCY_LAST_WRITE;
            default:
                throw new IllegalArgumentException("Missing StateConcurrency mapping to gRPC Concurrency enum");
        }
    }

    private RuntimeProto.StateOptions.StateConsistency getGrpcStateConsistency(StateOptions options) {
        switch (options.getConsistency()) {
            case EVENTUAL:
                return RuntimeProto.StateOptions.StateConsistency.CONSISTENCY_EVENTUAL;
            case STRONG:
                return RuntimeProto.StateOptions.StateConsistency.CONSISTENCY_STRONG;
            default:
                throw new IllegalArgumentException("Missing Consistency mapping to gRPC Consistency enum");
        }
    }

    /**
     * Execute a transaction.
     *
     * @param request Request to execute transaction.
     */
    @Override
    public void executeStateTransaction(ExecuteStateTransactionRequest request) {
        try {
            final String stateStoreName = request.getStateStoreName();
            final List<TransactionalStateOperation<?>> operations = request.getOperations();
            final Map<String, String> metadata = request.getMetadata();

            // 1. validate
            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }

            // 2. construct request object
            RuntimeProto.ExecuteStateTransactionRequest.Builder builder = RuntimeProto.ExecuteStateTransactionRequest.newBuilder();
            builder.setStoreName(stateStoreName);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            for (TransactionalStateOperation<?> operation : operations) {
                RuntimeProto.TransactionalStateOperation.Builder operationBuilder = RuntimeProto.TransactionalStateOperation.newBuilder();

                String operationType = operation.getOperation().toString().toLowerCase();
                operationBuilder.setOperationType(operationType);

                RuntimeProto.StateItem stateItem = buildStateRequest(operation.getRequest())
                        .build();
                operationBuilder.setRequest(stateItem);

                builder.addOperations(operationBuilder.build());
            }
            RuntimeProto.ExecuteStateTransactionRequest req = builder.build();

            // 3. invoke grpc api
            blockingStub.executeStateTransaction(req);
        } catch (Exception e) {
            logger.error("executeStateTransaction error ", e);
            throw new RuntimeClientException(e);
        }
    }

    /**
     * Retrieve a State based on their key.
     *
     * @param request The request to get state.
     * @param clazz   The Class of State needed as return.
     * @return The requested State.
     */
    @Override
    public <T> State<T> getState(GetStateRequest request, Class<T> clazz) {
        try {
            final String stateStoreName = request.getStoreName();
            final String key = request.getKey();
            final StateOptions options = request.getStateOptions();
            final Map<String, String> metadata = request.getMetadata();

            // 1. validate
            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }
            if ((key == null) || (key.trim().isEmpty())) {
                throw new IllegalArgumentException("Key cannot be null or empty.");
            }

            // 2. construct request object
            RuntimeProto.GetStateRequest.Builder builder = RuntimeProto.GetStateRequest.newBuilder()
                    .setStoreName(stateStoreName)
                    .setKey(key);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            if (options != null && options.getConsistency() != null) {
                RuntimeProto.StateOptions.StateConsistency consistency = getGrpcStateConsistency(options);
                builder.setConsistency(consistency);
            }
            RuntimeProto.GetStateRequest envelope = builder.build();

            // 3. invoke grpc api
            RuntimeProto.GetStateResponse resp = blockingStub.getState(envelope);

            // 4. parse result
            return parseGetStateResult(resp, key, options, clazz);
        } catch (Exception e) {
            logger.error("getState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    @Override
    public <T> List<State<T>> getBulkState(GetBulkStateRequest request, Class<T> clazz) {
        try {
            final String stateStoreName = request.getStoreName();
            final List<String> keys = request.getKeys();
            final int parallelism = request.getParallelism();
            final Map<String, String> metadata = request.getMetadata();

            // 1. validate
            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }
            if (keys == null || keys.isEmpty()) {
                throw new IllegalArgumentException("Key cannot be null or empty.");
            }
            if (parallelism < 0) {
                throw new IllegalArgumentException("Parallelism cannot be negative.");
            }

            // 2. construct request object
            RuntimeProto.GetBulkStateRequest.Builder builder = RuntimeProto.GetBulkStateRequest.newBuilder()
                    .setStoreName(stateStoreName)
                    .addAllKeys(keys)
                    .setParallelism(parallelism);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            RuntimeProto.GetBulkStateRequest envelope = builder.build();

            // 3. invoke grpc API
            RuntimeProto.GetBulkStateResponse resp = blockingStub.getBulkState(envelope);

            // 4. parse result
            List<RuntimeProto.BulkStateItem> itemsList = resp.getItemsList();
            List<State<T>> result = new ArrayList<>(itemsList.size());
            for (RuntimeProto.BulkStateItem itm : itemsList) {
                State<T> tState = parseGetStateResult(itm, clazz);
                result.add(tState);
            }
            return result;
        } catch (Exception e) {
            logger.error("getBulkState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    private <T> State<T> parseGetStateResult(
            RuntimeProto.GetStateResponse getStateResponse,
            String requestedKey,
            StateOptions stateOptions,
            Class<T> clazz) throws IOException {
        final ByteString payload = getStateResponse.getData();
        byte[] data = payload == null ? null : payload.toByteArray();
        T value = stateSerializer.deserialize(data, clazz);
        String etag = getStateResponse.getEtag();
        if (etag != null && etag.isEmpty()) {
            etag = null;
        }
        return new State<>(requestedKey, value, etag, getStateResponse.getMetadataMap(), stateOptions);
    }

    private <T> State<T> parseGetStateResult(
            RuntimeProto.BulkStateItem bulkStateItem,
            Class<T> clazz) throws IOException {
        final String key = bulkStateItem.getKey();
        final String error = bulkStateItem.getError();
        if (!Strings.isNullOrEmpty(error)) {
            return new State<>(key, error);
        }

        final ByteString payload = bulkStateItem.getData();
        byte[] data = payload == null ? null : payload.toByteArray();
        T value = stateSerializer.deserialize(data, clazz);
        String etag = bulkStateItem.getEtag();
        if (etag != null && etag.isEmpty()) {
            etag = null;
        }
        return new State<>(key, value, etag, bulkStateItem.getMetadataMap(), null);
    }
}
