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

    freeIndex := p.RequestEntry(v.ID)
    if freeIndex == -1 {
        return false
    }

    v.entered = true
    parkingTime := time.Duration(rand.Intn(3)+2) * time.Second
    time.Sleep(parkingTime)

    p.ExitVehicle(freeIndex)
    v.exited = true
    return true
}
