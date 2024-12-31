package module

import (
	"context"
	"errors"
	"plugin"
	"sync"

	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/utils"
)

const DefaultModuleDir = "modules"

type ModuleManager struct {
	options *Options
	modules map[string]Module
	mu      sync.RWMutex
}

func NewModuleManager(options ...Option) *ModuleManager {
	return &ModuleManager{
		modules: make(map[string]Module),
		options: newOptions(options...),
	}
}

func (mm *ModuleManager) RegisterModules(ctx context.Context, modules ...Module) error {
	for _, m := range modules {
		if err := mm.RegisterModule(ctx, m); err != nil {
			mm.options.Logger.Error("Failed to register module", logger.Error(err))
			return err
		}
	}
	return nil
}

func (mm *ModuleManager) RegisterModule(ctx context.Context, m Module) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	log := mm.options.Logger
	if err := m.Initialize(ctx, mm.options); err != nil {
		log.Error("Failed to initialize module", logger.Error(err))
		return err
	}

	// Ensure no duplicate module registration
	if _, exists := mm.modules[m.Id()]; exists {
		err := errors.New("duplicate module ID: " + m.Id())
		log.Error("Module registration failed", logger.Error(err))
		return err
	}

	mm.modules[m.Id()] = m
	log.Info("Module registered successfully",
		logger.String("id", m.Id()),
		logger.String("name", m.Name()),
		logger.String("version", m.Version()))
	return nil
}

func (mm *ModuleManager) LoadAllDynamicModules(ctx context.Context, modules []ModuleConfig) {
	log := mm.options.Logger
	log.Info("Loading dynamic modules")
	for _, module := range modules {
		if module.Enable {
			if module.Id == "" {
				log.Error("Module id is required", logger.String("name", module.Name))
				continue
			}
			err := mm.DynamicRegisterModuleById(ctx, module.Id)
			if err != nil {
				log.Error("Failed to load module", logger.String("id", module.Id), logger.String("name", module.Name), logger.Error(err))
			}
		}
	}
}

func (mm *ModuleManager) DynamicRegisterModuleById(ctx context.Context, moduleId string) error {
	path := DefaultModuleDir + "/" + moduleId + ".so"

	if !utils.FileExists(path) {
		return ErrModuleNotFound
	}

	module, err := mm.LoadModule(path)
	if err != nil {
		return err
	}

	return mm.RegisterModule(ctx, module)
}

func (mm *ModuleManager) LoadModule(modulePath string) (Module, error) {
	log := mm.options.Logger

	// Open the module file
	plug, err := plugin.Open(modulePath)
	if err != nil {
		log.Error("Failed to open module", logger.String("path", modulePath), logger.Error(err))
		return nil, err
	}

	// Look up the 'Module' symbol
	symModule, err := plug.Lookup("Module")
	if err != nil {
		log.Error("Failed to find 'Module' symbol in module", logger.String("path", modulePath), logger.Error(err))
		return nil, err
	}

	// Assert the type to module.Module
	loadedModule, ok := symModule.(Module)
	if !ok {
		err := ErrInvalidModuleType
		log.Error("Module type assertion failed", logger.String("path", modulePath), logger.String("expectedType", "module.Module"), logger.Error(err))
		return nil, err
	}
	return loadedModule, nil
}

func (mm *ModuleManager) GetModule(moduleId string) Module {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	return mm.modules[moduleId]
}

func (mm *ModuleManager) GetModules() map[string]Module {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	// Return a copy to avoid external modification
	copiedModules := make(map[string]Module)
	for id, module := range mm.modules {
		copiedModules[id] = module
	}
	return copiedModules
}

func (mm *ModuleManager) StartAllModules(ctx context.Context) error {
	log := mm.options.Logger
	log.Info("Starting all modules")
	for _, module := range mm.GetModules() {
		if err := module.Start(ctx); err != nil {
			log.Error("Failed to start module", logger.String("id", module.Id()), logger.Error(err))
			return err
		}
	}
	return nil
}
