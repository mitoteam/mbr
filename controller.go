package mbr

type Controller interface {
	With(middleware Middleware)
	Middlewares() MiddlewareList

	ParentControllers() []Controller
	SetParentControllers(parents []Controller)
}

type ControllerBase struct {
	middlewares       MiddlewareList
	parentControllers []Controller
}

// force interface implementation declaring empty pointer
var _ Controller = (*ControllerBase)(nil)

func (ctrl *ControllerBase) With(middleware Middleware) {
	ctrl.middlewares = append(ctrl.middlewares, middleware)
}

func (ctrl *ControllerBase) Middlewares() MiddlewareList {
	return ctrl.middlewares
}

func (ctrl *ControllerBase) ParentControllers() []Controller {
	return ctrl.parentControllers
}

func (ctrl *ControllerBase) SetParentControllers(parents []Controller) {
	ctrl.parentControllers = parents
}
