package httpapp

import (
	"github.com/RussellLuo/appx"
)

type App struct {
	*appx.App
}

func New(name string, instance appx.Instance) *App {
	return &App{App: appx.New(name, instance)}
}

func (a *App) MountOn(parent, pattern string) *App {
	m := MountOn(parent, pattern)
	a.App.Instance = appx.Standardize(m(a.App.Instance))

	a.App.Require(parent)
	return a
}

func (a *App) Require(names ...string) *App {
	a.App.Require(names...)
	return a
}
