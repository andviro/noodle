package noodle

// Default creates new middleware Chain with Recover middleware on top
func Default(mws ...Middleware) Chain {
	return New(Logger, Recover).Use(mws...)
}
