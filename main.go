package main

import (
    "AppFyne/src/models"
    "AppFyne/src/scenes"
    "fyne.io/fyne/v2/app"
)

func main() {
    a := app.New()
    w := a.NewWindow("Simulador de Estacionamiento")

    parking := models.NewParking(20)
    scene := scenes.NewParkingScene(w, parking)

    scene.Start()
    w.ShowAndRun()
}
