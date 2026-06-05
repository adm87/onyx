package file

import "path/filepath"

type FilePath string

func (f FilePath) String() string {
	return string(f)
}

func (f FilePath) Ext() FileExt {
	return FileExt(filepath.Ext(string(f)))
}

func (f FilePath) IsEmpty() bool {
	return f == ""
}

type FileExt string

func (e FileExt) String() string {
	return string(e)
}

func (e FileExt) IsEmpty() bool {
	return e == ""
}
