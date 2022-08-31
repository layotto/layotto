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
	s3ext "mosn.io/layotto/pkg/grpc/extension/s3"
	runtime_lock "mosn.io/layotto/pkg/runtime/lock"
)

type extensionComponentFactorys struct {
	locks []*runtime_lock.Factory
}

func WithExtensionGrpcAPI() Option {
	return WithGrpcAPI(
		s3ext.NewS3Server,
	)
}
