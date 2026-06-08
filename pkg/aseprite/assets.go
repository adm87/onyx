package aseprite

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/images"
)

const AdapterID engine.AdapterID = "aseprite_adapter"

type animation struct {
	meta      Meta
	frames    []AnimationFrame
	frameTags map[string]FrameTag
}

type AsepriteAssetAdapter struct {
	assets engine.Assets
	logger engine.Logger

	imageAssetAdapter *images.ImageAssetAdapter

	animations map[file.FilePath]*animation
}

func NewAsepriteAssetAdapter(assets engine.Assets, logger engine.Logger, imageAssetAdapter *images.ImageAssetAdapter) *AsepriteAssetAdapter {
	return &AsepriteAssetAdapter{
		assets:            assets,
		logger:            logger,
		imageAssetAdapter: imageAssetAdapter,
		animations:        make(map[file.FilePath]*animation),
	}
}

func (a *AsepriteAssetAdapter) ImportAsset(filesystem fs.FS, path file.FilePath, raw []byte) error {
	return nil // This asset adapter doesn't handle importing raw data.
}

func (a *AsepriteAssetAdapter) DeleteAsset(path file.FilePath) bool {
	return true // This asset adapter doesn't handle deleting raw data.
}

func (a *AsepriteAssetAdapter) SupportedExtensions() []file.FileExt {
	return nil // This asset adapter doesn't handle importing raw data.
}

func (a *AsepriteAssetAdapter) ImportAnimations(jsonPaths ...file.FilePath) error {
	for _, path := range jsonPaths {
		if _, exists := a.animations[path]; exists {
			continue
		}

		jsonData, exists := a.assets.GetData(path)
		if !exists {
			return fmt.Errorf("animation data not found for path: %s", path)
		}

		var data AsepriteJson
		if err := json.Unmarshal(jsonData, &data); err != nil {
			return err
		}

		if data.Meta.Image == "" {
			return fmt.Errorf("animation '%s' is missing image reference in meta", path)
		}

		imagePath := file.ResolvedPath(filepath.Dir(path.String()), data.Meta.Image.String())
		data.Meta.Image = imagePath

		if _, exists := a.imageAssetAdapter.GetImage(imagePath); !exists {
			return fmt.Errorf("referenced image '%s' for animation '%s' not found in image assets", imagePath, path)
		}

		anim := &animation{
			meta:      data.Meta,
			frames:    data.Frames,
			frameTags: make(map[string]FrameTag),
		}
		for _, tag := range data.Meta.FrameTags {
			anim.frameTags[tag.Name] = tag
		}

		a.logger.Debug("Imported Aseprite animation from '%s'", path)
		a.animations[path] = anim

		a.assets.Unload(path) // Unload the raw JSON data since it's no longer needed after importing.
	}
	return nil
}

func (a *AsepriteAssetAdapter) DeleteAnimations(jsonPaths ...file.FilePath) {
	for _, path := range jsonPaths {
		delete(a.animations, path)
	}
}
