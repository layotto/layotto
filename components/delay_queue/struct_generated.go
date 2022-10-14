// Code generated by github.com/layotto/protoc-gen-p6 .

// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package delay_queue

// DelayMessageRequest is the message to publish
type DelayMessageRequest struct {
	// Required. The name of the DelayQueue component
	ComponentName string `json:"component_name,omitempty"`
	// Required. The pubsub topic
	Topic string `json:"topic,omitempty"`
	// Required. The data which will be published to topic.
	Data []byte `json:"data,omitempty"`
	// The content type for the data (optional).
	DataContentType string `json:"data_content_type,omitempty"`
	// The length of time, in seconds, for which the delivery
	// of this messages is delayed.  Default: 0.
	DelayInSeconds int32 `json:"delay_in_seconds,omitempty"`
	// The metadata passing to pub components
	//
	// metadata property:
	// - key : the key of the message.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// DelayMessageResponse is the response
type DelayMessageResponse struct {
	// The message identifier
	MessageId string `json:"message_id,omitempty"`
}
