package smash

import (
	"context"
	"errors"
	"fmt"
	"plugin"
)

var (
	errModuleAlreadyLoaded = errors.New("module already loaded")
	errModuleNotLoaded     = errors.New("module not loaded")
)

type ModuleManager struct {
	bus       Bus
	scheduler *Scheduler
	modules   map[string]Module
}

func NewModuleManager(bus Bus, scheduler *Scheduler) (*ModuleManager, error) {
	return &ModuleManager{
		bus:       bus,
		scheduler: scheduler,
		modules:   map[string]Module{},
	}, nil
}

func (mm *ModuleManager) Load(modulePath string) (*Module, error) {
	p, err := plugin.Open(modulePath)

	if err != nil {
		return nil, err
	}

	registerSym, err := p.Lookup("Register")

	if err != nil {
		return nil, err
	}

	nameSym, err := p.Lookup("Name")

	if err != nil {
		return nil, err
	}

	registerFunc, ok := registerSym.(func(Bus, *Scheduler) error)

	if !ok {
		return nil, fmt.Errorf("cannot convert registerSym")
	}

	name, ok := nameSym.(*string)

	if !ok {
		return nil, fmt.Errorf("cannot convert nameSym")
	}

	module := Module{
		Name:     *name,
		Register: registerFunc,
	}

	if _, ok := mm.modules[module.Name]; ok {
		return nil, errModuleAlreadyLoaded
	}

	mm.modules[module.Name] = module

	_, err = mm.bus.Emit(context.Background(),
		EmitOptions{
			Message: Message{
				Type:   "module:register",
				Source: module.Name,
			},
		})

	if err != nil {
		return nil, err
	}

	return &module, nil
}

func (mm *ModuleManager) RegisterModules() error {
	for _, module := range mm.modules {
		if err := module.Register(mm.bus, mm.scheduler); err != nil {
			return err
		}
	}

	return nil
}

func (mm *ModuleManager) Unload(m Module) error {
	if _, ok := mm.modules[m.Name]; !ok {
		return errModuleNotLoaded
	}

	delete(mm.modules, m.Name)

	return nil
}

type Module struct {
	Name     string
	Register func(Bus, *Scheduler) error
}
