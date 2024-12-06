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

func (v *Vehicle) EnterParking(p *Parking) {
	v.mu.Lock()
	if v.entered || v.exited {
		v.mu.Unlock()
		return
	}
	v.mu.Unlock()

	// Intentar entrar al estacionamiento
	if p.RequestEntry(v.ID) {
		v.mu.Lock()
		v.entered = true
		v.mu.Unlock()

		// Tiempo aleatorio de estacionamiento
		parkingTime := time.Duration(rand.Intn(3)+2) * time.Second
		time.Sleep(parkingTime)

		// Salir del estacionamiento
		p.ExitVehicle(v.ID)

		v.mu.Lock()
		v.exited = true
		v.mu.Unlock()
	}
}