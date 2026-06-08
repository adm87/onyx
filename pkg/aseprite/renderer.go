package aseprite

import (
	"image"

	"github.com/adm87/onyx/pkg/engine"
	"github.com/adm87/onyx/pkg/engine/components/asset"
	"github.com/adm87/onyx/pkg/engine/components/rendering"
	"github.com/adm87/onyx/pkg/engine/components/transform"
	"github.com/adm87/onyx/pkg/engine/file"
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

const AsepriteRendererType rendering.RendererType = "aseprite_renderer"

type AsepriteRendererAdapter struct {
	imageAssetAdapter    *images.ImageAssetAdapter
	asepriteAssetAdapter *AsepriteAssetAdapter

	rendererTypes  []rendering.RendererType
	renderingTasks []engine.RenderTask
}

func NewAsepriteRendererAdapter(imageAssetAdapter *images.ImageAssetAdapter, asepriteAssetAdapter *AsepriteAssetAdapter) *AsepriteRendererAdapter {
	return &AsepriteRendererAdapter{
		imageAssetAdapter:    imageAssetAdapter,
		asepriteAssetAdapter: asepriteAssetAdapter,
		rendererTypes:        []rendering.RendererType{AsepriteRendererType},
		renderingTasks:       make([]engine.RenderTask, 0, 10),
	}
}

func (a *AsepriteRendererAdapter) SupportedRendererTypes() []rendering.RendererType {
	return a.rendererTypes
}

func (a *AsepriteRendererAdapter) GetRenderTasks(entry *donburi.Entry, layer int, zIndex int, viewport geom.AABB, viewMatrix ebiten.GeoM) []engine.RenderTask {
	a.renderingTasks = a.renderingTasks[:0]

	ref := asset.GetAssetReference(entry)
	if ref == asset.UnknownRef {
		return a.renderingTasks // Don't enqueue render tasks for entities without an image reference
	}

	anim, exists := a.asepriteAssetAdapter.animations[ref]
	if !exists {
		return a.renderingTasks // Don't enqueue render tasks for entities with an invalid animation reference
	}

	img, exists := a.imageAssetAdapter.GetImage(file.FilePath(anim.meta.Image))
	if !exists {
		return a.renderingTasks // Don't enqueue render tasks for entities with an invalid image reference
	}

	animName := GetAnimationName(entry)

	frameTag, exists := anim.frameTags[animName]
	if !exists {
		return a.renderingTasks // Don't enqueue render tasks for entities with an invalid animation tag reference
	}

	frames := anim.frames[frameTag.From : frameTag.To+1]
	if len(frames) == 0 {
		return a.renderingTasks // Don't enqueue render tasks for entities with an empty animation tag
	}

	frame := frames[GetFrameIndex(entry)]

	subImg := img.SubImage(
		image.Rect(
			frame.Frame.X,
			frame.Frame.Y,
			frame.Frame.X+frame.Frame.W,
			frame.Frame.Y+frame.Frame.H,
		),
	).(*ebiten.Image)

	anchor := rendering.GetAnchor(entry)

	scale := transform.GetScale(entry)

	aX := anchor.X * float64(frame.Frame.W) * scale.X
	aY := anchor.Y * float64(frame.Frame.H) * scale.Y

	matrix := transform.GetMatrix(entry)

	matrix.Translate(-aX, -aY)
	matrix.Concat(viewMatrix)

	a.renderingTasks = append(a.renderingTasks, engine.RenderTask{
		Render: func(screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM = matrix

			screen.DrawImage(subImg, opts)
			return nil
		},
		Layer:  layer,
		ZIndex: zIndex,
	})

	return a.renderingTasks
}
