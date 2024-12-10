package views

import (
	"AppFyne/src/models"
	"fmt"
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
	mainContainer  *fyne.Container  // Added this line
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

	// Create the grid container to hold car images
	ui.grid = container.NewWithoutLayout()

	// Load the parking lot image 
	parkingLotImg := canvas.NewImageFromFile("assets/estacionamiento.jpg")
	parkingLotImg.FillMode = canvas.ImageFillOriginal

	// Create the main container 
	ui.mainContainer = container.NewVBox(
		parkingLotImg,
		ui.grid,
		container.NewVBox(
			ui.statusLabel, 
			widget.NewLabel("Ocupación:"), 
			ui.progressBar, 
			ui.vehicleCounter,
		),
	)

	ui.window.SetContent(ui.mainContainer)
	ui.parking.AddObserver(ui)

	return ui
}


func (ui *UI) buildUI() {
	// Load the parking lot image
	parkingLotImg := canvas.NewImageFromFile("assets/estacionamiento.jpg")
	parkingLotImg.FillMode = canvas.ImageFillOriginal

	// Create the main container
	mainContainer := container.NewGridWithRows(2)
	mainContainer.Add(parkingLotImg)

	infoContainer := container.NewVBox(
		ui.statusLabel,
		widget.NewLabel("Ocupación:"),
		ui.progressBar,
		ui.vehicleCounter,
	)
	mainContainer.Add(infoContainer)

	ui.window.SetContent(mainContainer)
}

func (ui *UI) UpdateAvailableSpaces() {
	availableSpaces := ui.parking.Capacity - ui.parking.Occupied
	ui.statusLabel.SetText(fmt.Sprintf("Espacios disponibles: %d", availableSpaces))
	ui.progressBar.SetValue(float64(ui.parking.Occupied) / float64(ui.parking.Capacity))
	ui.window.Content().Refresh()
}

func CreateCompletedScene() fyne.CanvasObject {
	title := widget.NewLabel("Simulación de Estacionamiento Completada")
	title.TextStyle = fyne.TextStyle{Bold: true}

	description := widget.NewLabel("Todos los vehículos han completado su ciclo.")

	return container.NewVBox(
		title,
		description,
	)
}

func (ui *UI) animateCar(carImg *canvas.Image, fromX, fromY, toX, toY float32) {
	steps := 20 // Increase the number of steps for smoother animation
	stepX := (toX - fromX) / float32(steps)
	stepY := (toY - fromY) / float32(steps)

	// Channel to coordinate the animation
	done := make(chan struct{})

	go func() {
		for i := 1; i <= steps; i++ {
			currentX := fromX + stepX*float32(i)
			currentY := fromY + stepY*float32(i)

			// Use canvas.Refresh to update
			ui.window.Canvas().Refresh(carImg)
			carImg.Move(fyne.NewPos(currentX, currentY))

			time.Sleep(30 * time.Millisecond) // Adjust the delay for desired animation speed
		}
		close(done)
	}()

	<-done
}

func (ui *UI) ShowVehicle(id int) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	if _, exists := ui.carImages[id]; !exists {
		carImg := canvas.NewImageFromFile("assets/auto.png")
		carImg.FillMode = canvas.ImageFillContain

		spaceIndex := ui.parking.Occupied - 1
		x, y := ui.getGridPosition(spaceIndex)

		carImg.Move(fyne.NewPos(x, y))
		carImg.Resize(fyne.NewSize(70, 70)) 
		ui.grid.Add(carImg)
		ui.carImages[id] = carImg

		// Refresh the main container
		ui.mainContainer.Refresh()
	}
}

func (ui *UI) getGridPosition(index int) (float32, float32) {
	cols := 5
	if index == -1 {
		return -100, -100
	}
	x := float32(index%cols) * 80
	y := float32(index/cols) * 80
	return x, y
}
