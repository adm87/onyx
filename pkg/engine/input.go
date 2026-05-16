package engine

type Input interface {
}

type input struct {
	logger Logger
}

func NewInput(logger Logger) Input {
	return &input{
		logger: logger,
	}
}
