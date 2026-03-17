package connectors

import (
	"fmt"
	"sort"
)

type Registry struct {
	adapters map[string]Adapter
}

func NewRegistry(adapters ...Adapter) *Registry {
	items := map[string]Adapter{}
	for _, adapter := range adapters {
		items[adapter.Name()] = adapter
	}
	return &Registry{adapters: items}
}

func (r *Registry) Adapter(name string) (Adapter, error) {
	adapter, ok := r.adapters[name]
	if !ok {
		return nil, fmt.Errorf("connector adapter %s not found", name)
	}
	return adapter, nil
}

func (r *Registry) List() []string {
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
