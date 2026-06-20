package images

type ImagePlugin struct {
	assets *ImageAssets
}

func NewImagePlugin() *ImagePlugin {
	return &ImagePlugin{
		assets: NewImageAssets(),
	}
}

func (i *ImagePlugin) Assets() *ImageAssets {
	return i.assets
}
