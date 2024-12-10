package models

type Observer interface {
    UpdateAvailableSpaces()
}

type Subject interface {
    Register(observer Observer)
    Unregister(observer Observer)
    NotifyAll()
}