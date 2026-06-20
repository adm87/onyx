package ecs

import (
	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/adm87/onyx/pkg/plugins/ecs/camera"
	"github.com/adm87/onyx/pkg/plugins/ecs/image"
	"github.com/adm87/onyx/pkg/plugins/ecs/renderer"
	"github.com/adm87/onyx/pkg/plugins/ecs/transform"
	imageplugin "github.com/adm87/onyx/pkg/plugins/images"
	tiledplugin "github.com/adm87/onyx/pkg/plugins/tiled"
	"github.com/yohamta/donburi"
)

type ECSFactory struct {
	partitioner *ECSPartitioner

	imageAssets *imageplugin.ImageAssets
	tiledAssets *tiledplugin.TiledAssets

	imageRendererType uint64
	tiledRendererType uint64
}

func (f *ECSFactory) CreateCamera(world donburi.World, opts ...transform.TransformOption) *donburi.Entry {
	entry := transform.NewTransform(world, opts...)
	entry.AddComponent(camera.MainCamera)
	return entry
}

func (f *ECSFactory) CreateImage(world donburi.World, opts ...image.Option) *donburi.Entry {
	entry := image.NewImage(world, opts...)

	var bounds geom.AABB

	imgHandle := image.GetHandle(entry)
	frameIdx := image.GetFrame(entry)

	if img, exists := f.imageAssets.GetFrame(imgHandle, frameIdx); exists {
		anchor := image.GetAnchor(entry)

		width, height := img.Bounds().Dx(), img.Bounds().Dy()
		bounds.Min = geom.Vec2{
			X: -anchor.X * float64(width),
			Y: -anchor.Y * float64(height),
		}
		bounds.Max = geom.Vec2{
			X: bounds.Min.X + float64(width),
			Y: bounds.Min.Y + float64(height),
		}
	}

	// Note: By default entities are positioned at the world origin, so we can use the bounds as-is
	// without needing to offset them by the entity's position.
	idx, _ := f.partitioner.Add(entry.Entity(), bounds)

	transform.AddTransform(entry,
		transform.WithBounds(bounds.Min, bounds.Max),
		transform.WithIndex(idx),
	)

	renderer.AddRenderer(entry,
		renderer.WithRendererType(f.imageRendererType),
	)

	return entry
}
