package fundi

import (
	"github.com/goava/di"
)

type repository struct{}

func newRepository() *repository {
	return &repository{}
}

func provideRepository() di.Option {
	return di.Options(
		di.Provide(
			newRepository,
		),
	)
}
