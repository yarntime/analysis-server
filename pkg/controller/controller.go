package controller

type MTController struct {
}

func NewMTController(c *Config) *MTController {
	return &MTController{}
}

func (mtc *MTController) Run(stopCh chan struct{}) {

}
