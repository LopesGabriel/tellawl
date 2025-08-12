package controllers

type controller struct {
	version string
}

func NewController() *controller {
	return &controller{
		version: "v1",
	}
}
