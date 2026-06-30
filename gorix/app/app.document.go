package app

import (
	"fmt"

	"github.com/Gromosome/gorix/gorix/core/context"
)

func (a *App) connectDocuments(
	ctx *context.Context,
) error {
	configs, err := a.config.DocumentConfigs()
	if err != nil {
		return err
	}

	if len(configs) == 0 {
		return nil
	}

	for _, documentConfig := range configs {
		if err := a.documents.Connect(
			ctx,
			documentConfig,
		); err != nil {
			return fmt.Errorf(
				"gorix: document bootstrap failed: %w",
				err,
			)
		}
	}

	return nil
}
