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
	entrance      *ParkingEntrance
	availableCond *sync.Cond
}

func NewParking(capacity int) *Parking {
	p := &Parking{
		Capacity:   capacity,
		Occupied:   0,
		Vehicles:   make(map[int]bool),
		observers:  []Observer{},
		entrance:   NewParkingEntrance(),
	}
	p.availableCond = sync.NewCond(&p.mu)
	return p
}

func (p *Parking) RequestEntry(vehicleID int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	for p.Occupied >= p.Capacity {
		p.availableCond.Wait()
	}

	p.entrance.Enter()
	defer p.entrance.Exit()

	if !p.Vehicles[vehicleID] {
		p.Vehicles[vehicleID] = true
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
		p.notifyObservers()
		p.availableCond.Broadcast()
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


