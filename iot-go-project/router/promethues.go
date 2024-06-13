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
