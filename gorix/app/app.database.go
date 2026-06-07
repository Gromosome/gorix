package app

import (
	"fmt"

	"github.com/Gromosome/gorix/gorix/core/context"
)

func (a *App) connectDatabases(
	ctx *context.Context,
) error {
	configs, err := a.config.DatabaseConfigs()
	if err != nil {
		return err
	}

	for _, databaseConfig := range configs {
		if err := a.databases.Connect(
			ctx,
			databaseConfig,
		); err != nil {
			return fmt.Errorf(
				"gorix: database bootstrap failed: %w",
				err,
			)
		}
	}

	return nil
}
