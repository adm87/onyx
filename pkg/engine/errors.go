package engine

import "errors"

var (
	ErrAssetAdapterNotFound = errors.New("asset adapter not found")
	ErrAssetNotFound        = errors.New("asset not found")
)
