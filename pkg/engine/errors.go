package engine

type ErrAssetNotFound struct {
	Path string
}

func (e ErrAssetNotFound) Error() string {
	return "asset not found: " + e.Path
}

type ErrJsonUnmarshal struct {
	Path string
	Err  error
}

func (e ErrJsonUnmarshal) Error() string {
	return "failed to unmarshal json: " + e.Path + ": " + e.Err.Error()
}
