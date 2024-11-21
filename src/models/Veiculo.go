package models

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Vehicle struct {
    ID int
}

func NewVehicle(id int) *Vehicle {
    return &Vehicle{ID: id}
}

func (v *Vehicle) EnterParking(p *Parking) {
    maxAttempts := 10
    for attempt := 0; attempt < maxAttempts; attempt++ {
        if p.EnterVehicle(v.ID) {
            time.Sleep(time.Duration(rand.Intn(3)+3) * time.Second)
            p.ExitVehicle(v.ID)
            return
        }
        wait := time.Duration(math.Pow(2, float64(attempt))) * time.Second
        jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
        time.Sleep(wait + jitter)
    }
    fmt.Printf("Vehículo %d no pudo entrar al estacionamiento\n", v.ID)
}
