package scenes

import (
	"AppFyne/src/models"
	"AppFyne/src/views"
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

type ParkingScene struct {
    window   fyne.Window
    parking  *models.Parking
    ui       *views.UI
    done     chan struct{}
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
    ps.done = make(chan struct{})
    vehicleCount := 100
    
    go func() {
        var wg sync.WaitGroup
        for i := 1; i <= vehicleCount; i++ {
            wg.Add(1)
            vehicle := models.NewVehicle(i)
            
            go func(v *models.Vehicle) {
                defer wg.Done()
                ps.ui.StartVehicle(v.ID)
                v.EnterParking(ps.parking)
            }(vehicle)
            
            time.Sleep(300 * time.Millisecond) 
        }
        
        wg.Wait()
        close(ps.done)
    }()

    go func() {
        <-ps.done
        fmt.Println("Simulación completada")
    }()
}
