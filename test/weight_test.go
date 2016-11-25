package test

import (
	"../src/balance"
	"../src/core"
	"math"
	"math/rand"
	"net"
	"testing"
	"time"
)

type DummyContext struct{}

func (d DummyContext) String() string {
	return "123"
}

func (d DummyContext) Ip() net.IP {
	return make(net.IP, 1)
}

func (d DummyContext) Port() int {
	return 0
}

func TestWeightDistribution(t *testing.T) {
	rand.Seed(time.Now().Unix())
	balancer := &balance.WeightBalancer{}
	var context core.Context

	context = DummyContext{}

	backends := []core.Backend{
		{
			Weight: 20,
		},
		{
			Weight: 15,
		},
		{
			Weight: 25,
		},
		{
			Weight: 40,
		},
	}

	quantity := make(map[int]int)

	for _, backend := range backends {
		quantity[backend.Weight] = 0
	}

	n := 10000
	for try := 0; try < 100*n; try++ {
		backend, err := balancer.Elect(&context, backends)
		if err != nil {
			t.Fatal(err)
		}

		quantity[backend.Weight] += 1
	}

	for k, v := range quantity {
		if math.Abs(float64(v)/float64(n)-float64(k)) > 0.5 {
			t.Error(k, ":", float64(v)/float64(n))
		}

	}
}
