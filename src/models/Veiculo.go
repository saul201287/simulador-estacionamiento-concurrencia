package models

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"sync"
	"time"
)

type Vehicle struct {
    ID int
}

type UI struct {
	window    fyne.Window
	parking   *Parking
	grid      *fyne.Container
	carImages map[int]*canvas.Image
	mu        sync.Mutex
}

func NewVehicle(id int) *Vehicle {
    return &Vehicle{ID: id}
}

func (v *Vehicle) EnterParking(p *Parking) {
    if p.EnterVehicle(v.ID) {
        time.Sleep(3 * time.Second)
        p.ExitVehicle(v.ID)
    }
}
func (ui *UI) UpdateAvailableSpaces() {
    ui.mu.Lock()
    defer ui.mu.Unlock()

    // Actualizar la cuadrícula para mostrar los autos
    for id, img := range ui.carImages {
        if !ui.parking.Vehicles[id] { // Acceso al campo público Vehicles
            ui.grid.Remove(img) // Quitar imagen si el auto ha salido
            delete(ui.carImages, id)
        }
    }
    ui.window.Content().Refresh()
}
