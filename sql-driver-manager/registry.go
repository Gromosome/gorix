package sql_driver_manager

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

var adapters = struct {
	sync.RWMutex
	values map[string]Adapter
}{
	values: make(map[string]Adapter),
}

func Register(adapter Adapter) {
	if adapter == nil {
		panic("gorix sql-driver-manager: adapter cannot be nil")
	}

	name := strings.ToLower(strings.TrimSpace(adapter.Name()))
	sqlName := strings.TrimSpace(adapter.SQLDriverName())
	if name == "" || sqlName == "" {
		panic("gorix sql-driver-manager: adapter name and SQL driver name are required")
	}

	adapters.Lock()
	defer adapters.Unlock()

	if _, exists := adapters.values[name]; exists {
		panic("gorix sql-driver-manager: adapter already registered: " + name)
	}

	adapters.values[name] = adapter
}

func Lookup(name string) (Adapter, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))

	adapters.RLock()
	adapter, exists := adapters.values[normalized]
	adapters.RUnlock()

	if !exists {
		return nil, fmt.Errorf(
			"gorix sql-driver-manager: driver %q is not registered; blank-import its wrapper module",
			name,
		)
	}

	return adapter, nil
}

func RegisteredDrivers() []string {
	adapters.RLock()
	defer adapters.RUnlock()

	result := make([]string, 0, len(adapters.values))
	for name := range adapters.values {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
