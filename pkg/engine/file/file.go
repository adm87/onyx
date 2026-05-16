package file

import "path/filepath"

type Ext string

func (e Ext) IsEmpty() bool {
	return e == ""
}

type Path string

func (p Path) String() string {
	return string(p)
}

func (p Path) IsEmpty() bool {
	return p == ""
}

func (p Path) Ext() Ext {
	return Ext(filepath.Ext(p.String()))
}

func (p Path) Base() string {
	return filepath.Base(p.String())
}

func (p Path) Dir() string {
	return filepath.Dir(p.String())
}

func (p Path) Join(elem ...string) Path {
	return Path(filepath.Join(p.String(), filepath.Join(elem...)))
}
