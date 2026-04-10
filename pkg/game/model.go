package game

type Model interface {
	Title() string
	SetTitle(title string)

	Version() string
	SetVersion(version string)
}

type modelImpl struct {
	title   string
	version string
}

func (m *modelImpl) Title() string {
	return m.title
}

func (m *modelImpl) SetTitle(title string) {
	m.title = title
}

func (m *modelImpl) Version() string {
	return m.version
}

func (m *modelImpl) SetVersion(version string) {
	m.version = version
}
