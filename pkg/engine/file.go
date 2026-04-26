package engine

import "path/filepath"

type FilePath string

func (p FilePath) String() string {
	return string(p)
}

func (p FilePath) Type() FileType {
	return FileType(filepath.Ext(p.String())[1:])
}

type FileType string

func (t FileType) String() string {
	return string(t)
}

func (t FileType) IsEmpty() bool {
	return t == ""
}
