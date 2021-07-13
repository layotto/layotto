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

package actuator

import (
	"testing"

	"mosn.io/layotto/components/pkg/common"

	testify "github.com/stretchr/testify/assert"
)

func TestGetHealthInitOrSuccess(t *testing.T) {
	assert := testify.New(t)

	hi := GetRuntimeReadinessIndicator()
	v, _ := hi.Report()
	assert.Equal(v, common.INIT)
	hi.SetStarted()
	h, _ := hi.Report()
	assert.Equal(h, common.UP)

	hi = GetRuntimeLivenessIndicator()
	v, _ = hi.Report()
	assert.Equal(v, common.INIT)
	hi.SetStarted()
	h, _ = hi.Report()
	assert.Equal(h, common.UP)
}

func TestGetHealthErrorThenRecover(t *testing.T) {
	assert := testify.New(t)

	// readiness error
	hi := GetRuntimeReadinessIndicator()
	hi.SetUnhealthy("sub error")
	h, v := hi.Report()
	assert.Equal(h, common.DOWN)
	assert.Equal(v[reasonKey], "sub error")

	// readiness indicator can recover
	hi.SetHealthy("recover")
	h, v = hi.Report()
	assert.Equal(h, common.UP)
	assert.Equal(v[reasonKey], "recover")

	// liveness error
	hi = GetRuntimeLivenessIndicator()
	hi.SetUnhealthy("sub error")
	h, v = hi.Report()
	assert.Equal(h, common.DOWN)
	assert.Equal(v[reasonKey], "sub error")

}
