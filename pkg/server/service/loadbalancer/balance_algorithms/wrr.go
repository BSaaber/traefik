package balance_algorithms

import (
	"container/heap"
)

type strategyWRR struct {
	handlers    []*namedHandler
	curDeadline float64
}

func newStrategyWeightedRoundRobin() Strategy {
	return &strategyWRR{}
}

func (s *strategyWRR) nextServer(status map[string]struct{}) *namedHandler {
	var handler *namedHandler
	for {
		handler = heap.Pop(s).(*namedHandler)

		s.curDeadline = handler.deadline
		handler.deadline += 1 / handler.weight

		heap.Push(s, handler)
		if _, ok := status[handler.name]; ok {
			break
		}
	}
	return handler
}

func (s *strategyWRR) add(h *namedHandler) {
	h.deadline = s.curDeadline + 1/h.weight
	heap.Push(s, h)
}

func (s *strategyWRR) setUp(string, bool) {}

func (s *strategyWRR) name() string {
	return "wrr"
}

func (s *strategyWRR) len() int {
	return len(s.handlers)
}

// Len implements heap.Interface/sort.Interface.
func (s *strategyWRR) Len() int { return s.len() }

// Less implements heap.Interface/sort.Interface.
func (s *strategyWRR) Less(i, j int) bool {
	return s.handlers[i].deadline < s.handlers[j].deadline
}

func (s *strategyWRR) Swap(i, j int) {
	s.handlers[i], s.handlers[j] = s.handlers[j], s.handlers[i]
}

func (s *strategyWRR) Push(x interface{}) {
	h, ok := x.(*namedHandler)
	if !ok {
		return
	}

	s.handlers = append(s.handlers, h)
}

func (s *strategyWRR) Pop() interface{} {
	h := s.handlers[len(s.handlers)-1]
	s.handlers = s.handlers[0 : len(s.handlers)-1]
	return h
}