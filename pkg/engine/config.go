package engine

type Config struct {
	Title      string
	Width      int
	Height     int
	Fullscreen bool
}

func NewConfig() *Config {
	return &Config{
		Title:      "Untitled",
		Width:      800,
		Height:     600,
		Fullscreen: false,
	}
}

func (c *Config) WithTitle(title string) *Config {
	c.Title = title
	return c
}

func (c *Config) WithWidth(width int) *Config {
	c.Width = width
	return c
}

func (c *Config) WithHeight(height int) *Config {
	c.Height = height
	return c
}

func (c *Config) WithFullscreen(fullscreen bool) *Config {
	c.Fullscreen = fullscreen
	return c
}
