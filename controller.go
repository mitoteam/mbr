package mbr

type Controller interface {
	BasePath() string
}

type ControllerBase struct {
	parent Controller
}

// force interface implementation declaring empty pointer
var _ Controller = (*ControllerBase)(nil)

func (c *ControllerBase) BasePath() string {
	return "/"
}
