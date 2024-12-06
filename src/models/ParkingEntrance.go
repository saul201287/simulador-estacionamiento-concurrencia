package models

import (
	"sync"
)

type ParkingEntrance struct {
	mu        sync.Mutex
	occupied  bool
	enterCond *sync.Cond
}

func NewParkingEntrance() *ParkingEntrance {
	pe := &ParkingEntrance{}
	pe.enterCond = sync.NewCond(&pe.mu)
	return pe
}

func (pe *ParkingEntrance) Enter() {
	pe.mu.Lock()
	defer pe.mu.Unlock()

	for pe.occupied {
		pe.enterCond.Wait()
	}
	pe.occupied = true
}

func (pe *ParkingEntrance) Exit() {
	pe.mu.Lock()
	pe.occupied = false
	pe.enterCond.Broadcast()
	pe.mu.Unlock()
}
