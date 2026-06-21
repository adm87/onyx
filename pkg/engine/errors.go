package engine

type ErrAssetNotFound struct {
	Path string
}

func (e ErrAssetNotFound) Error() string {
	return "asset not found: " + e.Path
}
