package balance_algorithms

import (
	crand "crypto/rand"
	"math/rand/v2"
)

type rnd interface {
	IntN(n int) int
}

type strategyPowerOfTwoChoices struct {
	healthy   []*namedHandler
	unhealthy []*namedHandler
	rand      rnd
}

func newStrategyPowerOfTwoRandomChoises() Strategy {
	return &strategyPowerOfTwoChoices{
		rand: newRand(),
	}
}

func (s *strategyPowerOfTwoChoices) nextServer(map[string]struct{}) *namedHandler {
	if len(s.healthy) == 1 {
		return s.healthy[0]
	}

	n1, n2 := s.rand.IntN(len(s.healthy)), s.rand.IntN(len(s.healthy)-1)
	if n2 >= n1 {
		n2 = (n2 + 1) % len(s.healthy)
	}

	h1, h2 := s.healthy[n1], s.healthy[n2]
	if h2.inflight.Load() < h1.inflight.Load() {
		return h2
	}

	return h1
}

func (s *strategyPowerOfTwoChoices) add(h *namedHandler) {
	s.healthy = append(s.healthy, h)
}

func (s *strategyPowerOfTwoChoices) setUp(name string, up bool) {
	if up {
		var healthy *namedHandler
		healthy, s.unhealthy = deleteAndPop(s.unhealthy, name)
		s.healthy = append(s.healthy, healthy)
		return
	}

	var unhealthy *namedHandler
	unhealthy, s.healthy = deleteAndPop(s.healthy, name)
	s.unhealthy = append(s.unhealthy, unhealthy)
}

func (s *strategyPowerOfTwoChoices) name() string {
	return "p2c"
}

func (s *strategyPowerOfTwoChoices) len() int {
	return len(s.healthy) + len(s.unhealthy)
}

func newRand() *rand.Rand {
	var seed [16]byte
	_, err := crand.Read(seed[:])
	if err != nil {
		panic(err)
	}
	var seed1, seed2 uint64
	for i := 0; i < 16; i += 8 {
		seed1 = seed1<<8 + uint64(seed[i])
		seed2 = seed2<<8 + uint64(seed[i+1])
	}
	return rand.New(rand.NewPCG(seed1, seed2))
}

func deleteAndPop(handlers []*namedHandler, name string) (deleted *namedHandler, remaining []*namedHandler) {
	for i, h := range handlers {
		if h.name == name {
			handlers[i], handlers[len(handlers)-1] = handlers[len(handlers)-1], handlers[i]
			deleted = handlers[len(handlers)-1]
			remaining = handlers[:len(handlers)-1]
			return
		}
	}

	panic("impossible to go here")
}