package mbr

type Controller interface {
}

type ControllerBase struct {
}

// force interface implementation declaring empty pointer
var _ Controller = (*ControllerBase)(nil)
