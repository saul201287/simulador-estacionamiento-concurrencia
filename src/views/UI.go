package views

import (
	"AppFyne/src/models"
	"fmt"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)


type UI struct {
	window         fyne.Window
	parking        *models.Parking
	grid           *fyne.Container
	carImages      map[int]*canvas.Image
	statusLabel    *widget.Label
	progressBar    *widget.ProgressBar
	vehicleCounter *widget.Label
	mu             sync.Mutex
}


func NewUI(window fyne.Window, parking *models.Parking) *UI {
	ui := &UI{
		window:         window,
		parking:        parking,
		carImages:      make(map[int]*canvas.Image),
		statusLabel:    widget.NewLabel("Espacios disponibles: 20"),
		progressBar:    widget.NewProgressBar(),
		vehicleCounter: widget.NewLabel("Vehículos en espera: 0"),
	}
	
	ui.parking.AddObserver(ui)
	ui.buildUI()
	return ui
}

func (ui *UI) buildUI() {
	rows, cols := 4, 5
	ui.grid = container.NewGridWithColumns(cols)

	for i := 0; i < rows*cols; i++ {
		rect := canvas.NewRectangle(color.NRGBA{G: 128, A: 255})
		rect.SetMinSize(fyne.NewSize(80, 80))
		ui.grid.Add(rect)
	}

	infoContainer := container.NewVBox(
		ui.statusLabel,
		widget.NewLabel("Ocupación:"),
		ui.progressBar,
		ui.vehicleCounter,
	)

	mainContainer := container.NewHSplit(ui.grid, infoContainer)
	ui.window.SetContent(mainContainer)
}

func (ui *UI) UpdateAvailableSpaces() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	availableSpaces := ui.parking.Capacity - ui.parking.Occupied
	ui.statusLabel.SetText(fmt.Sprintf("Espacios disponibles: %d", availableSpaces))

	ui.progressBar.SetValue(float64(ui.parking.Occupied) / float64(ui.parking.Capacity))

	for i, child := range ui.grid.Objects {
		rect, ok := child.(*canvas.Rectangle)
		if ok {
			if ui.isSpaceOccupied(i) {
				rect.FillColor = color.NRGBA{R: 255, A: 255} 
			} else {
				rect.FillColor = color.NRGBA{G: 128, A: 255} 
			}
			rect.Refresh()
		}
	}

	waitingVehicles := 100 - ui.parking.Occupied 
	ui.vehicleCounter.SetText(fmt.Sprintf("Vehículos en espera: %d", waitingVehicles))

	ui.window.Content().Refresh()
}

func (ui *UI) isSpaceOccupied(spaceIndex int) bool {
	return spaceIndex < ui.parking.Occupied
}

func (ui *UI) ShowVehicle(id int) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	if _, exists := ui.carImages[id]; !exists {
		carImg := canvas.NewImageFromFile("assets/auto.png")
		carImg.FillMode = canvas.ImageFillContain
		ui.grid.Add(carImg)
		ui.carImages[id] = carImg
		ui.window.Content().Refresh()
	}
}

func (ui *UI) StartVehicle(id int) {
	go func() {
		ui.ShowVehicle(id)
		time.Sleep(3 * time.Second)
		ui.parking.ExitVehicle(id)
	}()
}