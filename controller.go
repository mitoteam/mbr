package mbr

type Controller interface {
	With(middleware Middleware)
	Middlewares() MiddlewareList
}

type ControllerBase struct {
	middlewares MiddlewareList
}

// force interface implementation declaring empty pointer
var _ Controller = (*ControllerBase)(nil)

func (ctrl *ControllerBase) With(middleware Middleware) {
	ctrl.middlewares = append(ctrl.middlewares, middleware)
}

func (ctrl *ControllerBase) Middlewares() MiddlewareList {
	return ctrl.middlewares
}
