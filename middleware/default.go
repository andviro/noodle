package middleware

import (
	"log"

	"gopkg.in/andviro/noodle.v2"
)

// Default is a convenience function creating new noodle.Chain with Logger, Recover and LocalStore middlewares
func Default(mws ...noodle.Middleware) noodle.Chain {
	return noodle.New(RealIP, Logger, Recover(func(err error) {
		log.Printf("Error in handler: %#v\n", err)
	}), LocalStore).Use(mws...)
}
