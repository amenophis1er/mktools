package plugin

import (
	"context"
	"github.com/spf13/cobra"
)

// Plugin represents a tool that can be executed by mktools
type Plugin interface {
	// Name returns the unique identifier of the plugin
	Name() string

	// Description returns a short description of what the plugin does
	Description() string

	// Execute runs the plugin with the given context, command, and arguments
	Execute(ctx context.Context, cmd *cobra.Command, args []string) error

	// AddFlags adds plugin-specific flags to the command
	AddFlags(cmd *cobra.Command)
}

// Registry manages registered plugins
type Registry struct {
	plugins map[string]Plugin
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

// Register adds a plugin to the registry
func (r *Registry) Register(p Plugin) {
	r.plugins[p.Name()] = p
}

// Get retrieves a plugin by name
func (r *Registry) Get(name string) (Plugin, bool) {
	p, ok := r.plugins[name]
	return p, ok
}

// List returns all registered plugins
func (r *Registry) List() []Plugin {
	plugins := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}
