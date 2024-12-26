package mbr

type Controller interface {
	WithMiddlewareList

	ParentControllers() []Controller
	SetParentControllers(parents []Controller)
}

type ControllerBase struct {
	WithMiddlewareListBase

	parentControllers []Controller
}

// force interface implementation declaring empty pointer
var _ Controller = (*ControllerBase)(nil)

func (ctrl *ControllerBase) ParentControllers() []Controller {
	return ctrl.parentControllers
}

func (ctrl *ControllerBase) SetParentControllers(parents []Controller) {
	ctrl.parentControllers = parents
}
