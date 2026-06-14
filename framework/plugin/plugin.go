package plugin

import (
	"gpma/framework/data"
	"gpma/framework/errors"
)

type PluginResult struct {
	ASTModified bool
	MapModified bool
	Errors      []*errors.Error
	FatalError  error
}

func NewPluginResult() *PluginResult {
	return &PluginResult{
		ASTModified: false,
		MapModified: false,
		Errors:      nil,
		FatalError:  nil,
	}
}

func (r *PluginResult) WithASTModified() *PluginResult {
	r.ASTModified = true
	return r
}

func (r *PluginResult) WithMapModified() *PluginResult {
	r.MapModified = true
	return r
}

func (r *PluginResult) WithErrors(errs []*errors.Error) *PluginResult {
	r.Errors = append(r.Errors, errs...)
	return r
}

func (r *PluginResult) WithFatalError(err error) *PluginResult {
	r.FatalError = err
	return r
}

func (r *PluginResult) HasModifications() bool {
	return r.ASTModified || r.MapModified
}

func (r *PluginResult) HasErrors() bool {
	return len(r.Errors) > 0 || r.FatalError != nil
}

func (r *PluginResult) HasCriticalErrors() bool {
	for _, err := range r.Errors {
		if err.Severity == errors.SeverityCritical {
			return true
		}
	}
	return r.FatalError != nil
}

type Plugin interface {
	Name() string
	Process(astData *data.ASTData) *PluginResult
	Configure(params map[string]interface{}) error
}

type ConfigurablePlugin struct{}

func (cp *ConfigurablePlugin) Configure(params map[string]interface{}) error {
	return nil
}
