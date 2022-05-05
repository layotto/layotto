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

package http

import (
	"context"
	"encoding/json"

	"github.com/valyala/fasthttp"
	"mosn.io/api"
	mosnhttp "mosn.io/mosn/pkg/protocol/http"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/variable"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/actuator"
)

type DispatchFilter struct {
	handler api.StreamReceiverFilterHandler
}

func (dis *DispatchFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	dis.handler = handler
}

func (dis *DispatchFilter) OnDestroy() {}

func (dis *DispatchFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	// 1. log
	log.DefaultLogger.Debugf("[actuator] receive actuator pkt")
	path, err := variable.GetString(ctx, types.VarHttpRequestPath)
	if err != nil {
		dis.write404()
		return api.StreamFilterStop
	}
	log.DefaultLogger.Debugf("[actuator] path: %v", path)
	// 2. validate path
	resolver := NewPathResolver(path)
	// http path must be /actuator/{endpoint_name}/{params}
	// So we can return 404 directly if it does not start with "actuator"
	if resolver.Next() != "actuator" {
		// illegal
		dis.write404()
		return api.StreamFilterStop
	}
	act := actuator.GetDefault()
	if act == nil {
		dis.write404()
		return api.StreamFilterStop
	}
	// 3. dispatch endpoint
	epName := resolver.Next()
	endpoint, ok := act.GetEndpoint(epName)
	if !ok {
		// illegal
		dis.write404()
		return api.StreamFilterStop
	}
	json, err := endpoint.Handle(ctx, resolver)
	// 4. write result
	var code int
	if err != nil {
		code = HttpUnavailableCode
	} else {
		code = HttpSuccessCode
	}
	dis.writeJsonResult(json, code)
	return api.StreamFilterStop
}

func (dis *DispatchFilter) write404() {
	dis.writeJsonResult(nil, HttpNotFoundCode)
}

func (dis *DispatchFilter) writeJsonResult(jsonObject map[string]interface{}, code int) {
	if code == 0 {
		code = HttpSuccessCode
	}
	// 0. marshal
	var byteSlice []byte
	if jsonObject != nil {
		var err error
		byteSlice, err = json.Marshal(jsonObject)
		if err != nil {
			log.DefaultLogger.Errorf("[actuator][dispatch_filter]error when marshal result:%v", err)
			code = HttpUnavailableCode
		}
	}
	// 1. header
	fastHttpHeader := &fasthttp.ResponseHeader{}
	rspHeader := mosnhttp.ResponseHeader{
		ResponseHeader: fastHttpHeader,
	}
	rspHeader.Set("Content-Type", "application/json")
	rspHeader.SetStatusCode(code)
	// 2. body
	data := buffer.NewIoBufferBytes(byteSlice)
	// 3. write response
	dis.handler.SendDirectResponse(rspHeader, data, nil)
}
