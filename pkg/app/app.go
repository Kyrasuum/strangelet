package app

import ()

type App struct {
	StartApp func()
	Priv     interface{}
	View     interface{}
}
