package game

type Model struct {
	title   string
	version string
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) Title() string {
	return m.title
}

func (m *Model) SetTitle(title string) {
	m.title = title
}

func (m *Model) Version() string {
	return m.version
}

func (m *Model) SetVersion(version string) {
	m.version = version
}
