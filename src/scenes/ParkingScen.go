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
    // Simular ingreso de vehículos cada cierto tiempo
    go func() {
        for i := 1; i <= 10; i++ { // Generar 10 vehículos
            vehicle := models.NewVehicle(i)
            ps.ui.StartVehicle(vehicle.ID)
            vehicle.EnterParking(ps.parking)
            time.Sleep(2 * time.Second) // Tiempo entre ingresos
        }
    }()
}
