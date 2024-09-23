package process

type Process struct {
	pId       int
	user      string
	cdm       string
	cpuUseage float32
	memory    float32

	// TODO: change the type to upTime to time
	upTime string
}

func NewProcess() *Process {
	// TODO: fill the required fields
	return &Process{}
}

