package engine

type InputBindingID uint64

type InputBinding interface {
	SetActive(active bool)
	Poll() error
	ID() InputBindingID
}

type InputBindings struct {
}

type Input struct {
	logger   *Logger
	bindings map[InputBindingID]InputBinding
}

func NewInput(logger *Logger) *Input {
	return &Input{
		logger:   logger,
		bindings: make(map[InputBindingID]InputBinding),
	}
}

func (s *Input) EnableBinding(id InputBindingID) {
	if binding, exists := s.bindings[id]; exists {
		binding.SetActive(true)
		return
	}
	s.logger.Warn("Attempted to enable non-existent input binding with ID %d", id)
}

func (s *Input) DisableBinding(id InputBindingID) {
	if binding, exists := s.bindings[id]; exists {
		binding.SetActive(false)
		return
	}
	s.logger.Warn("Attempted to disable non-existent input binding with ID %d", id)
}

func (s *Input) Poll() error {
	for _, binding := range s.bindings {
		if err := binding.Poll(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Input) Bind(bindings ...InputBinding) {
	for _, binding := range bindings {
		if _, exists := s.bindings[binding.ID()]; exists {
			s.logger.Warn("Input binding with ID %d already exists, overwriting", binding.ID())
		}
		s.bindings[binding.ID()] = binding
		s.logger.Debug("Bound input with ID", "id", binding.ID())
	}
}
