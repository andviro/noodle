package middleware

import (
	"github.com/andviro/noodle"
)

// Default is a convenience function creating new noodle.Chain with Logger, Recover and LocalStore middlewares
func Default(mws ...noodle.Middleware) noodle.Chain {
	return noodle.New(RealIP, Logger, Recover, LocalStore).Use(mws...)
}
