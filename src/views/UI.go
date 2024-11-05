package views

import (
	"AppFyne/src/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"sync"
	"time"
)

type UI struct {
	window    fyne.Window
	parking   *models.Parking
	grid      *fyne.Container
	carImages map[int]*canvas.Image
	mu        sync.Mutex
}

func NewUI(window fyne.Window, parking *models.Parking) *UI {
	ui := &UI{
		window:    window,
		parking:   parking,
		carImages: make(map[int]*canvas.Image),
	}
	ui.parking.AddObserver(ui)
	ui.buildUI()
	return ui
}

func (ui *UI) buildUI() {
	// Crear una cuadrícula para el estacionamiento
	rows, cols := 4, 5 // Ajusta el tamaño del estacionamiento
	ui.grid = container.NewGridWithColumns(cols)

	for i := 0; i < rows*cols; i++ {
		rect := canvas.NewRectangle(color.NRGBA{G: 128, A: 255})
		rect.SetMinSize(fyne.NewSize(80, 80))
		ui.grid.Add(rect)
	}

	mainContainer := container.NewBorder(nil, nil, nil, nil, ui.grid)
	ui.window.SetContent(mainContainer)
}

func (ui *UI) UpdateAvailableSpaces() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	// Actualizar la cuadrícula para mostrar los autos
	for id, img := range ui.carImages {
		if !ui.parking.Vehicles[id] {
			ui.grid.Remove(img) // Quitar imagen si el auto ha salido
			delete(ui.carImages, id)
		}
	}
	ui.window.Content().Refresh()
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
