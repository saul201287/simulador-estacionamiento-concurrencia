package models

import (
    "time"
    "math/rand"
)

type Vehicle struct {
    ID int
}

func NewVehicle(id int) *Vehicle {
    return &Vehicle{ID: id}
}

func (v *Vehicle) EnterParking(p *Parking) {
    // Intentar ingresar al estacionamiento
    if p.EnterVehicle(v.ID) {
        time.Sleep(time.Duration(rand.Intn(3)+3) * time.Second) // Estacionar entre 3 y 5 segundos
        p.ExitVehicle(v.ID)
    } else {
        // Si no pudo entrar, intenta nuevamente después de un tiempo
        time.Sleep(2 * time.Second) // Intentar cada 2 segundos
        v.EnterParking(p)
    }
}
