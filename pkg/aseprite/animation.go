package aseprite

import (
	"github.com/adm87/onyx/pkg/engine/components/asset"
	"github.com/yohamta/donburi"
)

type AnimationInfo struct {
	AnimationName string
	FrameIndex    int

	animationTime float64
}

var Animation = donburi.NewComponentType[AnimationInfo]()

func NewAnimation(ecs donburi.World, animationName string) *donburi.Entry {
	return AddAnimation(ecs.Entry(
		ecs.Create(
			Animation,
		),
	), animationName)
}

func AddAnimation(entry *donburi.Entry, animationName string) *donburi.Entry {
	donburi.Add(entry, Animation, &AnimationInfo{
		AnimationName: animationName,
		FrameIndex:    0,
	})
	return entry
}

func GetAnimationName(entry *donburi.Entry) string {
	if !entry.HasComponent(Animation) {
		return ""
	}
	return Animation.Get(entry).AnimationName
}

func SetAnimationName(entry *donburi.Entry, name string) {
	if !entry.HasComponent(Animation) {
		return
	}

	animInfo := Animation.Get(entry)
	if animInfo.AnimationName == name {
		return
	}

	animInfo.AnimationName = name
	animInfo.FrameIndex = 0
	Animation.Set(entry, animInfo)
}

func GetFrameIndex(entry *donburi.Entry) int {
	if !entry.HasComponent(Animation) {
		return 0
	}
	return Animation.Get(entry).FrameIndex
}

func SetFrameIndex(entry *donburi.Entry, index int) {
	if !entry.HasComponent(Animation) {
		return
	}
	animInfo := Animation.Get(entry)
	animInfo.FrameIndex = index
	animInfo.animationTime = 0
	Animation.Set(entry, animInfo)
}

func GetAnimationFrame(entry *donburi.Entry, adapter *AsepriteAssetAdapter) AnimationFrame {
	ref := asset.GetAssetReference(entry)
	if ref == asset.UnknownRef {
		return AnimationFrame{}
	}

	anim, exists := adapter.animations[ref]
	if !exists {
		return AnimationFrame{}
	}

	animName := GetAnimationName(entry)
	frameIndex := GetFrameIndex(entry)

	if frameIndex >= len(anim.frames) {
		return AnimationFrame{}
	}

	frameTag, exists := anim.frameTags[animName]
	if !exists {
		return AnimationFrame{}
	}

	frames := anim.frames[frameTag.From : frameTag.To+1]
	if len(frames) == 0 {
		return AnimationFrame{}
	}

	return frames[frameIndex]
}

func UpdateAnimation(entry *donburi.Entry, adapter *AsepriteAssetAdapter, dt float64) {
	animInfo := Animation.Get(entry)

	ref := asset.GetAssetReference(entry)
	if ref == asset.UnknownRef {
		return
	}

	anim, exists := adapter.animations[ref]
	if !exists {
		return
	}

	frameTag, exists := anim.frameTags[animInfo.AnimationName]
	if !exists {
		return
	}

	frames := anim.frames[frameTag.From : frameTag.To+1]
	if len(frames) == 0 {
		return
	}

	animInfo.animationTime += dt
	frameDuration := float64(frames[animInfo.FrameIndex].Duration) / 1000.0

	for animInfo.animationTime >= frameDuration {
		animInfo.animationTime -= frameDuration
		animInfo.FrameIndex++
		if animInfo.FrameIndex >= len(frames) {
			animInfo.FrameIndex = 0
		}
		frameDuration = float64(frames[animInfo.FrameIndex].Duration) / 1000.0
	}

	Animation.Set(entry, animInfo)
}
