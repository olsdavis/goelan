package plugin

// Plugin struct represents an extension for goelan which can handle events,
// create new commands, etc. In short: customize the server thanks to an API
// for the used language.
type Plugin struct {
	Language Language
}

// Enable enables the current plugin.
func (plugin *Plugin) Enable() {
	// TODO: implement.
}

// Disable disables the current plugin and unloads the used resources.
func (plugin *Plugin) Disable() {
	// TODO: implement.
}

type PluginManager struct {
	plugins []Plugin
}

// LoadPlugins loads all the plugins from the given folder path.
// If some plugins have already been loaded, they are disabled and unloaded.
func (pm *PluginManager) LoadPlugins(folder string) (bool, error) {
	// there are loaded plugins
	if len(pm.plugins) > 0 {
		for _, plugin := range pm.plugins {
			plugin.Disable()
		}
	}

	// TODO: define how plugins can be recognized.

	return true, nil
}
