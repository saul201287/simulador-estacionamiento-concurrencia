package scenes

import (
    "AppFyne/src/models"
    "AppFyne/src/views"
    "time"
    "fyne.io/fyne/v2"
)

type ParkingScene struct {
    window fyne.Window
    parking *models.Parking
    ui      *views.UI
}

func NewParkingScene(window fyne.Window, parking *models.Parking) *ParkingScene {
    ui := views.NewUI(window, parking)
    return &ParkingScene{
        window:  window,
        parking: parking,
        ui:      ui,
    }
}

func (ps *ParkingScene) Start() {
    go func() {
        for i := 1; i <= 100; i++ {
            vehicle := models.NewVehicle(i)
            ps.ui.StartVehicle(vehicle.ID)
            go vehicle.EnterParking(ps.parking) // Cada vehículo en su propia goroutine
            time.Sleep(2 * time.Second) // Los vehículos llegan con intervalos de 2 segundos
        }
    }()
}
