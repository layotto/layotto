/*
 * Copyright 2021 Layotto Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package io.mosn.layotto.v1.callback.component.pubsub;

import java.util.Collection;
import java.util.HashSet;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class PubSubClientRegistryImpl implements PubSubRegistry {

    public final Map<String, PubSub> pubSubClients = new ConcurrentHashMap<>();

    @Override
    public void registerPubSubCallback(String pubsubName, PubSub callback) {
        if (pubsubName == null) {
            throw new IllegalArgumentException("pubSubName shouldn't be null");
        }
        if (callback == null) {
            throw new IllegalArgumentException("callback shouldn't be null");
        }
        if (pubSubClients.putIfAbsent(pubsubName, callback) != null) {
            throw new IllegalArgumentException("Pub/sub callback with name " + pubsubName + " already exists!");
        }
    }

    @Override
    public PubSub getCallbackByPubSubName(String pubsubName) {
        if (pubsubName == null) {
            throw new IllegalArgumentException("pubsubName shouldn't be null");
        }
        final PubSub pubSub = pubSubClients.get(pubsubName);
        if (pubSub == null) {
            throw new IllegalArgumentException("Cannot find pubsub callback by name " + pubsubName);
        }
        return pubSub;
    }

    @Override
    public Collection<PubSub> getAllPubSubCallbacks() {
        return new HashSet<>(pubSubClients.values());
    }
}
