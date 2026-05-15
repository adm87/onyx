package engine

type Assets interface {
}

type assets struct {
	logger Logger
}

func NewAssets(logger Logger) Assets {
	return &assets{
		logger: logger,
	}
}
