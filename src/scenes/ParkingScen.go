package scenes

import (
    "AppFyne/src/models"
    "AppFyne/src/views"
    "log"
    "sync"
    "time"

    "fyne.io/fyne/v2"
)

type ParkingScene struct {
    window  fyne.Window
    parking *models.Parking
    ui      *views.UI
    done    chan struct{}
}

func NewParkingScene(window fyne.Window, parking *models.Parking) *ParkingScene {
    ui := views.NewUI(window, parking)
    return &ParkingScene{
        window:  window,
        parking: parking,
        ui:      ui,
        done:    make(chan struct{}),
    }
}

func (ps *ParkingScene) Start() {
    vehicleCount := 100
    vehicles := make([]*models.Vehicle, vehicleCount)

    go func() {
        var wg sync.WaitGroup
        for i := 1; i <= vehicleCount; i++ {
            wg.Add(1)
            go func(vehicleID int) {
                defer wg.Done()
                vehicle := models.NewVehicle(vehicleID)
                vehicles[vehicleID-1] = vehicle
                for {
                    freeIndex := ps.parking.RequestEntry(vehicleID)
                    if freeIndex != -1 {
                        vehicle.EnterParking(ps.parking)
                        log.Printf("Carro %d estacionado", vehicleID)
                        break
                    }
                    time.Sleep(100 * time.Millisecond) // Reintentar después de un tiempo
                }
            }(i)
        }
        wg.Wait()
        close(ps.done) // Marca que la simulación ha terminado
    }()

    // Configura la escena de finalización cuando termine la simulación
    go func() {
        <-ps.done
        ps.window.SetContent(views.CreateCompletedScene())
    }()
}