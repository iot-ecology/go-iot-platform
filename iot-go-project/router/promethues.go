/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package router

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
)

func Mem() {
}

func NewRelicConfig() *newrelic.Application {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("管理后台"),
		newrelic.ConfigLicense("mit"),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if nil != err {
		fmt.Printf("New Relic initialization failed: %v", err)
		os.Exit(1)
	}

	return app
}
