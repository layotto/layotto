// Code generated by github.com/seeflood/protoc-gen-p6 .

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

package runtime

import (
	email "mosn.io/layotto/components/email"
	phone "mosn.io/layotto/components/phone"
	email1 "mosn.io/layotto/pkg/grpc/email"
	s3 "mosn.io/layotto/pkg/grpc/extension/s3"
	phone1 "mosn.io/layotto/pkg/grpc/phone"
)

type extensionComponentFactorys struct {
	// "mosn.io/layotto/spec/proto/extension/v1/email"
	// email.
	email []*email.Factory

	// "mosn.io/layotto/spec/proto/extension/v1/phone"
	// phone.
	phone []*phone.Factory
}

func WithEmailServiceFactory(email ...*email.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.email = append(o.services.email, email...)
	}
}

func WithPhoneCallServiceFactory(phone ...*phone.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.phone = append(o.services.phone, phone...)
	}
}

func WithExtensionGrpcAPI() Option {
	return WithGrpcAPI(
		s3.NewS3Server,
		email1.NewAPI,
		phone1.NewAPI,
	)
}
