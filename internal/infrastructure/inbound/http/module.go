package http //nolint:revive // intentional package name to group HTTP handlers

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			newCategoryHandler,
		),
	)
}
