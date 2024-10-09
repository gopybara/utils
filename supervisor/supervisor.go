package supervisor

import (
	"sync"
)

type GoroutineSupervisor struct {
	wg *sync.WaitGroup
}

func NewGoroutineSupervisor() *GoroutineSupervisor {
	return &GoroutineSupervisor{wg: new(sync.WaitGroup)}
}

func (s *GoroutineSupervisor) AddProc() {
	s.wg.Add(1)
}

func (s *GoroutineSupervisor) DoneProc() {
	s.wg.Done()
}

func (s *GoroutineSupervisor) WaitForComplete() {
	s.wg.Wait()
}
