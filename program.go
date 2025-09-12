package main

import (
	"oriongo/internal/origongo"
)

type Program struct {
	app *origongo.OrionGo
}

func NewProgram() *Program {
	return &Program{}
}

func (p *Program) StartApp() *Program {
	app := createDefaultApp()
	program := &Program{
		app: app,
	}

	return program
}

func createDefaultApp() *origongo.OrionGo {
	app := origongo.CreateDefaultApp()

	return app
}

func (p *Program) Run() *error {
	if p.app == nil {
		p.app = createDefaultApp()
	}
	p.app.Run()

	return nil
}
