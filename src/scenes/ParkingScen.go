package scenes

import (
	"AppFyne/src/models"
	"AppFyne/src/views"
	"fyne.io/fyne/v2"
	"sync"
	"time"
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

	go func() {
		var wg sync.WaitGroup
		semaphore := make(chan struct{}, 10)

		for i := 1; i <= vehicleCount; i++ {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(vehicleID int) {
				defer wg.Done()
				vehicle := models.NewVehicle(vehicleID)
				ps.ui.StartVehicle(vehicleID)
				vehicle.EnterParking(ps.parking)
				<-semaphore
			}(i)

			time.Sleep(200 * time.Millisecond)
		}

		wg.Wait()
		close(ps.done)
	}()

	go func() {
		<-ps.done
		ps.window.SetContent(views.CreateCompletedScene())
	}()
}
