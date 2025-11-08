package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CounterVec struct {
	counters *prometheus.CounterVec
}

func NewCounterVec(metricsName, help string, labels []string, opts ...OptsFunc) *CounterVec {
	opt := &mOpts{
		name: metricsName,
		help: help,
	}
	for _, optsFunc := range opts {
		optsFunc(opt)
	}
	counterOpt := opt.GetCounterOpts()
	cc := prometheus.NewCounterVec(counterOpt, labels)
	prometheus.MustRegister(cc)

	return &CounterVec{
		counters: cc,
	}
}

func (self *CounterVec) Inc(labels ...string) {
	self.counters.WithLabelValues(labels...).Inc()
}

func (self *CounterVec) Add(count float64, labels ...string) {
	self.counters.WithLabelValues(labels...).Add(count)
}

func (self *CounterVec) Delete(labels ...string) {
	self.counters.DeleteLabelValues(labels...)
}

type Counter struct {
	counter prometheus.Counter
}

func NewCounter(metricsName, help string, opts ...OptsFunc) *Counter {
	opt := &mOpts{
		name: metricsName,
		help: help,
	}
	for _, optsFunc := range opts {
		optsFunc(opt)
	}
	counterOpt := opt.GetCounterOpts()
	cc := prometheus.NewCounter(counterOpt)
	prometheus.MustRegister(cc)

	return &Counter{
		counter: cc,
	}
}

func (self *Counter) Inc() {
	self.counter.Inc()
}

func (self *Counter) Add(count float64) {
	self.counter.Add(count)
}
