package models

import (
    "sync"
)

type Parking struct {
    Capacity     int
    Occupied     int
    mu           sync.Mutex
    Vehicles     map[int]bool
    observers    []Observer
    parkingSlots chan struct{}
}

func NewParking(capacity int) *Parking {
    return &Parking{
        Capacity:     capacity,
        Occupied:     0,
        Vehicles:     make(map[int]bool),
        parkingSlots: make(chan struct{}, capacity),
    }
}

func (p *Parking) EnterVehicle(id int) bool {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.Occupied < p.Capacity {
        p.parkingSlots <- struct{}{} 
        p.Vehicles[id] = true
        p.Occupied++
        p.notifyObservers()
        return true
    }

    return false
}


func (p *Parking) ExitVehicle(id int) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.Vehicles[id] {
        delete(p.Vehicles, id)
        p.Occupied--
        <-p.parkingSlots 
        p.notifyObservers()
    }
}

func (p *Parking) notifyObservers() {
    for _, observer := range p.observers {
        observer.UpdateAvailableSpaces()
    }
}

func (p *Parking) AddObserver(o Observer) {
    p.observers = append(p.observers, o)
}
