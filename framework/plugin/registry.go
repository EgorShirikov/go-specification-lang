package plugin

import (
	"fmt"
	"sync"
)

type PluginFactory func() Plugin

var (
	registry = make(map[string]PluginFactory)
	mu       sync.RWMutex
)

func Register(name string, factory PluginFactory) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = factory
}

func Get(name string) (Plugin, error) {
	return GetConfigured(name, nil)
}

func GetConfigured(name string, params map[string]interface{}) (Plugin, error) {
	mu.RLock()
	defer mu.RUnlock()

	factory, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("plugin %q not found", name)
	}

	plugin := factory()

	if len(params) > 0 {
		if err := plugin.Configure(params); err != nil {
			return nil, fmt.Errorf("failed to configure plugin %q: %w", name, err)
		}
	}

	return plugin, nil
}

func GetAll(names []string) ([]Plugin, error) {
	return GetAllConfigured(names, nil)
}

func GetAllConfigured(names []string, configs map[string]map[string]interface{}) ([]Plugin, error) {
	plugins := make([]Plugin, 0, len(names))

	for _, name := range names {
		var pluginConfig map[string]interface{}
		if configs != nil {
			pluginConfig = configs[name]
		}

		plugin, err := GetConfigured(name, pluginConfig)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}

func ConfigurePlugin(plugin Plugin, params map[string]interface{}) error {
	return plugin.Configure(params)
}

func ConfigureAll(plugins []Plugin, configs map[string]map[string]interface{}) error {
	for _, plugin := range plugins {
		if configs != nil {
			if pluginConfig, exists := configs[plugin.Name()]; exists {
				if err := plugin.Configure(pluginConfig); err != nil {
					return fmt.Errorf("failed to configure plugin %q: %w", plugin.Name(), err)
				}
			}
		}
	}
	return nil
}

func List() []string {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
