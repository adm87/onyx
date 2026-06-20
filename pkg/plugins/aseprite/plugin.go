package aseprite

import (
	"encoding/json"
	"image"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/plugins/images"
)

type AsepritePlugin struct {
	logger     engine.Logger
	animations map[uint64]*AnimationData
}

func NewAsepritePlugin(logger engine.Logger) *AsepritePlugin {
	return &AsepritePlugin{
		logger:     logger,
		animations: make(map[uint64]*AnimationData),
	}
}

func (m *AsepritePlugin) BuildAnimations(imageAssets *images.ImageAssets, imgHandle uint64, data []byte) *AnimationData {
	var animations *AnimationData

	err := json.Unmarshal(data, &animations)
	if err != nil {
		m.logger.Error("failed to parse Aseprite JSON data: %v", err)
		return nil
	}

	animations.Meta.Clips = make(map[string]FrameTag, len(animations.Meta.FrameTags))
	for i := range animations.Meta.FrameTags {
		tag := animations.Meta.FrameTags[i]
		animations.Meta.Clips[tag.Name] = tag
	}

	rects := make([]image.Rectangle, len(animations.Frames))
	for i, frame := range animations.Frames {
		rects[i] = image.Rect(
			frame.Frame.X,
			frame.Frame.Y,
			frame.Frame.X+frame.Frame.W,
			frame.Frame.Y+frame.Frame.H,
		)
	}

	imageAssets.ExtractFrames(imgHandle, rects)
	m.animations[imgHandle] = animations

	return animations
}

func (m *AsepritePlugin) DeleteAnimations(handle uint64) {
	delete(m.animations, handle)
}
