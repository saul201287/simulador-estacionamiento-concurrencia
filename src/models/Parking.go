package models

import (
    "sync"
)

type Parking struct {
    Capacity      int
    Occupied      int
    mu            sync.Mutex
    Vehicles      map[int]bool
    observers     []Observer
    availableCond *sync.Cond
}

func NewParking(capacity int) *Parking {
    p := &Parking{
        Capacity:   capacity,
        Occupied:   0,
        Vehicles:   make(map[int]bool),
        observers:  []Observer{},
    }
    p.availableCond = sync.NewCond(&p.mu)
    return p
}

func (p *Parking) RequestEntry(vehicleID int) int {
    p.mu.Lock()
    defer p.mu.Unlock()

    // Espera hasta que haya un espacio disponible
    for p.Occupied >= p.Capacity {
        p.availableCond.Wait()
    }

    // Encuentra un espacio vacío
    for i := 0; i < p.Capacity; i++ {
        if !p.Vehicles[i] { // Busca el primer espacio libre
            p.Vehicles[i] = true
            p.Occupied++
            p.notifyObservers()
            return i
        }
    }

    return -1 // Debe manejarse correctamente en caso de error inesperado
}

func (p *Parking) ExitVehicle(spaceIndex int) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if spaceIndex < 0 || spaceIndex >= p.Capacity {
        return // Indice fuera de rango
    }

    if p.Vehicles[spaceIndex] {
        delete(p.Vehicles, spaceIndex)
        p.Occupied--
        p.notifyObservers()
        p.availableCond.Broadcast() // Notifica que hay un espacio libre
    }
}

func (p *Parking) notifyObservers() {
    for _, observer := range p.observers {
        observer.UpdateAvailableSpaces()
    }
}

func (p *Parking) AddObserver(o Observer) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.observers = append(p.observers, o)
}