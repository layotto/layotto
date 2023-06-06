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

package nacos

import (
	"mosn.io/pkg/log"
)

// An adapter to implement log.LoggerInterface in agollo package.
type DefaultLogger struct {
	logger log.ErrorLogger
}

func NewDefaultLogger(logger log.ErrorLogger) *DefaultLogger {
	return &DefaultLogger{
		logger: logger,
	}
}
func (d *DefaultLogger) Debugf(format string, params ...interface{}) {
	d.logger.Debugf(format, params)
}

func (d *DefaultLogger) Infof(format string, params ...interface{}) {
	d.logger.Infof(format, params)
}

func (d *DefaultLogger) Warnf(format string, params ...interface{}) {
	d.logger.Warnf(format, params)
}

func (d *DefaultLogger) Errorf(format string, params ...interface{}) {
	d.logger.Errorf(format, params)
}

func (d *DefaultLogger) Debug(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	format := buildForamat(v)
	d.logger.Debugf(format, v)
}

func buildForamat(v []interface{}) string {
	return "%v"
}

func (d *DefaultLogger) Info(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	format := buildForamat(v)
	d.logger.Infof(format, v)
}

func (d *DefaultLogger) Warn(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	format := buildForamat(v)
	d.logger.Warnf(format, v)
}

func (d *DefaultLogger) Error(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	format := buildForamat(v)
	d.logger.Errorf(format, v)
}
