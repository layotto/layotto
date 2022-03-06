/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package dapr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	mock_pubsub "mosn.io/layotto/pkg/mock/components/pubsub"
	"testing"
	"time"
)

func TestDaprGrpcAPIPublishEvent(t *testing.T) {
	t.Run("invalid pubsub name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		api := NewDaprServer("", nil,
			nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil,
			nil, nil, nil, nil, nil, nil)
		_, err := api.PublishEvent(context.Background(), &dapr_v1pb.PublishEventRequest{})
		assert.Equal(t, "rpc error: code = InvalidArgument desc = pubsub name is empty", err.Error())
	})

	t.Run("invalid topic", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		api := NewDaprServer("", nil, nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil,
			nil, nil, nil, nil, nil, nil)
		req := &dapr_v1pb.PublishEventRequest{
			PubsubName: "mock",
		}
		_, err := api.PublishEvent(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = topic is empty in pubsub mock", err.Error())
	})

	t.Run("component not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		api := NewDaprServer("", nil, nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil,
			nil, nil, nil, nil, nil, nil)
		req := &dapr_v1pb.PublishEventRequest{
			PubsubName: "abc",
			Topic:      "abc",
		}
		_, err := api.PublishEvent(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = pubsub abc not found", err.Error())
	})

	t.Run("publish success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		mockPubSub.EXPECT().Publish(gomock.Any()).Return(nil)
		mockPubSub.EXPECT().Features().Return(nil)
		api := NewDaprServer("", nil, nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil,
			nil, nil, nil, nil, nil, nil)
		req := &dapr_v1pb.PublishEventRequest{
			PubsubName: "mock",
			Topic:      "abc",
		}
		_, err := api.PublishEvent(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("publish net error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		mockPubSub.EXPECT().Publish(gomock.Any()).Return(fmt.Errorf("net error"))
		mockPubSub.EXPECT().Features().Return(nil)
		api := NewDaprServer("", nil, nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil,
			nil, nil, nil, nil, nil, nil)
		req := &dapr_v1pb.PublishEventRequest{
			PubsubName: "mock",
			Topic:      "abc",
		}
		_, err := api.PublishEvent(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = error when publish to topic abc in pubsub mock: net error", err.Error())
	})
}

func TestMosnRuntime_publishMessageGRPC(t *testing.T) {
	t.Run("drop it when publishing an expired message", func(t *testing.T) {
		cloudEvent := map[string]interface{}{
			pubsub.IDField:              "id",
			pubsub.SourceField:          "source",
			pubsub.DataContentTypeField: "content-type",
			pubsub.TypeField:            "type",
			pubsub.SpecVersionField:     "v1.0.0",
			pubsub.DataBase64Field:      "bGF5b3R0bw==",
			pubsub.ExpirationField:      time.Now().Add(-time.Minute).Format(time.RFC3339),
		}

		data, err := json.Marshal(cloudEvent)
		assert.Nil(t, err)

		msg := &pubsub.NewMessage{
			Data:     data,
			Topic:    "layotto",
			Metadata: make(map[string]string),
		}
		a := NewDaprServer("", nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil)

		var apiForTest = a.(*daprGrpcAPI)
		apiForTest.json = jsoniter.ConfigFastest
		// execute
		err = apiForTest.publishMessageGRPC(context.Background(), msg)
		// validate
		assert.Nil(t, err)
	})
}
