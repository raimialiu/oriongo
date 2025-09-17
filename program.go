package main

import (
	"oriongo/internal/infrastructure"
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
	app := origongo.CreateDefaultApp().AddControllers().AddDbContext(infrastructure.ConnectionConfig{
		AutoConnect: true,
		Host:        "127.0.0.1",
		Port:        "3306",
		Database:    "DDumper",
		Username:    "root",
		Password:    "DVorak@23000",
	},
	)

	return app
}

func (p *Program) Run() *error {
	if p.app == nil {
		p.app = createDefaultApp()
	}
	p.app.Run()

	return nil
}
