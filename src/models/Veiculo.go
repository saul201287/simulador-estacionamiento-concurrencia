package models

import (
	"math/rand"
	"sync"
	"time"
)

type Vehicle struct {
	ID      int
	entered bool
	exited  bool
	mu      sync.Mutex
}

func NewVehicle(id int) *Vehicle {
	return &Vehicle{
		ID: id,
	}
}

func (v *Vehicle) EnterParking(p *Parking) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.entered || v.exited {
		return false
	}

	if p.EnterVehicle(v.ID) {
		v.entered = true

		parkingTime := time.Duration(rand.Intn(4)+2) * time.Second
		time.Sleep(parkingTime)

		p.ExitVehicle(v.ID)

		v.exited = true
		return true
	}

	return false
}
