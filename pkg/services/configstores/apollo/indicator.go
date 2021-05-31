package apollo

import (
	"github.com/layotto/layotto/pkg/actuator/health"
	"sync"
)

const (
	reasonKey = "reason"
)

var (
	readinessIndicator *healthIndicator
	livenessIndicator  *healthIndicator
)

func init() {
	readinessIndicator = newHealthIndicator()
	livenessIndicator = newHealthIndicator()
	health.AddReadinessIndicator("apollo", readinessIndicator)
	health.AddLivenessIndicator("apollo", livenessIndicator)
}

func newHealthIndicator() *healthIndicator {
	return &healthIndicator{
		started: false,
		isErr:   false,
	}
}

func getReadinessIndicator() *healthIndicator {
	return readinessIndicator
}

func getLivenessIndicator() *healthIndicator {
	return livenessIndicator
}

type healthIndicator struct {
	mu sync.Mutex

	started   bool
	isErr     bool
	errReason string
}

func (idc *healthIndicator) Report() health.Health {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	if idc.isErr {
		h := health.NewHealth(health.DOWN)
		h.Details[reasonKey] = idc.errReason
		return h
	}
	if idc.started {
		return health.NewHealth(health.UP)
	}

	return health.NewHealth(health.INIT)
}

func (idc *healthIndicator) reportError(reason string) {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	if idc.isErr {
		return
	}
	idc.isErr = true
	idc.errReason = reason
}

func (idc *healthIndicator) setStarted() {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	idc.started = true
}
