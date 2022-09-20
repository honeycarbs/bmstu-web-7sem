package context

import (
	"neatly/internal/config"
	"neatly/pkg/logging"
	"sync"
)

type AppContext struct {
	Config *config.Config
}

var instance *AppContext
var once sync.Once

func GetInstance() *AppContext {
	once.Do(func() {
		logging.GetLogger().Info("Initializing application context")
		instance = &AppContext{
			Config: config.GetConfig(),
		}
	})

	return instance
}
