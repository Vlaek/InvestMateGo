package app

import (
	"fmt"
	"os"
	"strings"

	"invest-mate/internal/assets"
	"invest-mate/internal/users"
	"invest-mate/pkg/logger"
)

const (
	ModuleAssets = "assets"
	ModuleUsers  = "users"

	// Префикс модуля
	ConfigEnablePrefix = "ENABLE_MODULE_"
	ConfigEnableAll    = "ENABLE_ALL_MODULES"
)

// Все доступные модули приложения (в порядке инициализации)
var availableModules = []string{
	ModuleAssets,
	ModuleUsers,
}

// Конфигурация модуля
type ModuleConfig struct {
	Enabled  bool
	Name     string
	Priority int
}

// Регистрация всех модулей приложения
func (app *App) RegisterModules() {
	logger.InfoLog("Registering application modules...")

	// Получение конфигурации модулей
	moduleConfigs := app.getModuleConfigs()

	// Регистрация модулей
	for _, config := range moduleConfigs {
		app.registerModuleFromConfig(config)
	}

	logger.InfoLog("Total modules registered: %d", len(app.modules))
}

// Создание и регистрация модуля по конфигурации
func (app *App) registerModuleFromConfig(config ModuleConfig) {
	// Фабрика модулей
	var module AppModule

	switch config.Name {
	case ModuleAssets:
		module = &assets.ModuleWrapper{}
	case ModuleUsers:
		module = &users.ModuleWrapper{}
	default:
		logger.ErrorLog("Unknown module type: %s", config.Name)
		return
	}

	app.registerModuleWithConfig(config, module)
}

// Регистрация модуля с проверкой конфигурации
func (app *App) registerModuleWithConfig(config ModuleConfig, module AppModule) {
	if !config.Enabled {
		logger.InfoLog("Module '%s' is disabled, skipping", config.Name)
		return
	}

	app.registerModule(config.Name, module)
}

// Регистрация модуля
func (app *App) registerModule(name string, module AppModule) {
	if app.modules == nil {
		app.modules = make(map[string]AppModule)
	}

	// Проверка дублирования
	if _, exists := app.modules[name]; exists {
		logger.InfoLog("Module '%s' is already registered, skipping", name)
		return
	}

	app.modules[name] = module
	logger.InfoLog("Module '%s' registered", name)
}

// Получение конфигурации для всех доступных модулей
func (app *App) getModuleConfigs() []ModuleConfig {
	configs := make([]ModuleConfig, 0, len(availableModules))

	// Проверяем, включены ли все модули сразу
	enableAll := strings.ToLower(os.Getenv(ConfigEnableAll)) == "true"

	// Создаем конфигурацию для каждого модуля
	for i, moduleName := range availableModules {
		config := ModuleConfig{
			Name:     moduleName,
			Priority: i,
			Enabled:  app.isModuleEnabled(moduleName, enableAll),
		}
		configs = append(configs, config)
	}

	return configs
}

// Функция проверяет, включен ли модуль
func (app *App) isModuleEnabled(moduleName string, enableAll bool) bool {
	if enableAll {
		return true
	}

	configKey := ConfigEnablePrefix + strings.ToUpper(moduleName)

	defaultEnabled := true

	envValue := os.Getenv(configKey)
	if envValue == "" {
		return defaultEnabled
	}

	// Парсим значение переменной окружения
	switch strings.ToLower(envValue) {
	case "true", "1", "yes", "on", "enabled":
		return true
	case "false", "0", "no", "off", "disabled":
		return false
	default:
		logger.InfoLog("Invalid value for %s: '%s', using default: %v",
			configKey, envValue, defaultEnabled)
		return defaultEnabled
	}
}

// Инициализация всех зарегистрированных модулей
func (app *App) InitializeModules() error {
	if len(app.modules) == 0 {
		logger.ErrorLog("No modules registered - application will run without functionality")
		return nil
	}

	logger.InfoLog("Initializing %d registered modules...", len(app.modules))

	modulesToInitialize := app.getModulesByPriority()

	successCount := 0
	failureCount := 0

	for _, name := range modulesToInitialize {
		module := app.modules[name]
		logger.InfoLog("Initializing module: %s", name)

		if err := module.Initialize(app.DB, app.Config); err != nil {
			errorMsg := fmt.Sprintf("Module %s initialization failed: %v", name, err)

			if app.Config.Env == "production" {
				if app.isCriticalModule(name) {
					logger.ErrorLog("%s - stopping application", errorMsg)
					return fmt.Errorf("critical module %s failed to initialize: %w", name, err)
				}
				logger.ErrorLog("%s - continuing without module", errorMsg)
			} else {
				logger.ErrorLog("%s - continuing in %s mode", errorMsg, app.Config.Env)
			}
			failureCount++
		} else {
			logger.InfoLog("✅ Module %s initialized successfully", name)
			successCount++
		}
	}

	logger.InfoLog("Modules initialization completed: %d successful, %d failed",
		successCount, failureCount)

	if successCount == 0 && len(app.modules) > 0 {
		logger.ErrorLog("All modules failed to initialize")
	}

	return nil
}

// Возвращение имён модулей в порядке приоритета
func (app *App) getModulesByPriority() []string {
	configMap := make(map[string]ModuleConfig)
	for _, config := range app.getModuleConfigs() {
		configMap[config.Name] = config
	}

	// Сортируем модули по приоритету
	modules := make([]string, 0, len(app.modules))
	for name := range app.modules {
		modules = append(modules, name)
	}

	// TODO: Кринж))) Простая сортировка пузырьком по приоритету
	for i := 0; i < len(modules)-1; i++ {
		for j := i + 1; j < len(modules); j++ {
			prioI := configMap[modules[i]].Priority
			prioJ := configMap[modules[j]].Priority
			if prioI > prioJ {
				modules[i], modules[j] = modules[j], modules[i]
			}
		}
	}

	return modules
}

// Функция проверяет, является ли модуль критическим
func (app *App) isCriticalModule(name string) bool {
	criticalModules := map[string]bool{
		ModuleUsers:  true,
		ModuleAssets: true,
	}

	return criticalModules[name]
}

// Закрытие ресурсов всех модулей
func (app *App) CloseModules() {
	if len(app.modules) == 0 {
		return
	}

	logger.InfoLog("Closing %d modules...", len(app.modules))

	modulesToClose := app.getModulesByPriority()
	for i := len(modulesToClose) - 1; i >= 0; i-- {
		name := modulesToClose[i]
		module := app.modules[name]

		logger.InfoLog("Closing module: %s", name)
		if err := module.Close(); err != nil {
			logger.ErrorLog("Failed to close module %s: %v", name, err)
		} else {
			logger.InfoLog("Module %s closed successfully", name)
		}
	}

	logger.InfoLog("All modules closed")
}
