package monitor

import "context"

var apmContext context.Context

func GetApmContext() context.Context {
	return apmContext
}

func SetApmContext(ctx context.Context) {
	apmContext = ctx
}
