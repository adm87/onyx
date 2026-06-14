package aseprite

import (
	"encoding/json"
	"fmt"
	"image"
	"time"

	"github.com/adm87/onyx/pkg/engine/assert"
	"github.com/adm87/onyx/pkg/images"
	"github.com/yohamta/donburi"
)

type AsepriteModule struct {
	imageModule *images.ImageModule
	animations  map[uint64]*AnimationData
}

func NewAsepriteModule(imageModule *images.ImageModule) *AsepriteModule {
	return &AsepriteModule{
		imageModule: imageModule,
		animations:  make(map[uint64]*AnimationData),
	}
}

func (m *AsepriteModule) BuildAnimations(imgHandle uint64, data []byte) *AnimationData {
	var animations *AnimationData

	err := json.Unmarshal(data, &animations)
	assert.Nil(err, fmt.Sprintf("failed to parse Aseprite JSON data: %v", err))

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

	m.imageModule.ExtractFrames(imgHandle, rects)
	m.animations[imgHandle] = animations

	return animations
}

func (m *AsepriteModule) DeleteAnimations(handle uint64) {
	delete(m.animations, handle)
}

func (m *AsepriteModule) UpdateAnimation(entry *donburi.Entry, dt time.Duration) {
	if !IsPlaying(entry) {
		return
	}

	library, exists := m.getLibrary(entry)
	if !exists {
		return
	}

	clip := GetClip(entry)

	tag, exists := library.Meta.Clips[clip]
	if !exists {
		return
	}

	animator := GetAnimator(entry)
	if animator == nil {
		return
	}

	elapsed := animator.time + dt

	frame := GetAnimationFrame(entry)
	frameIndex := tag.From + frame

	duration := time.Duration(library.Frames[frameIndex].Duration) * time.Millisecond
	if duration <= 0 || elapsed < duration {
		animator.time = elapsed
		SetAnimator(entry, animator)
		return
	}

	nextFrame := frame
	frameCount := tag.To - tag.From + 1

	var completed bool
	for elapsed >= duration {
		elapsed -= duration

		nextFrame, completed = m.getNextFrame(nextFrame, animator, frameCount)
		if completed {
			break
		}

		frameIndex = tag.From + nextFrame
		duration = time.Duration(library.Frames[frameIndex].Duration) * time.Millisecond
	}

	if nextFrame != frame {
		images.SetFrame(entry, frameIndex)
		SetAnimationFrame(entry, nextFrame)
	}
	if animator.time != elapsed {
		animator.time = elapsed
		SetAnimator(entry, animator)
	}
	if completed {
		SetAnimationState(entry, AnimationStateStopped)
	}
}

func (m *AsepriteModule) getNextFrame(current int, animator *AnimatorInfo, frameCount int) (int, bool) {
	current += animator.direction

	loopComplete := current >= frameCount || current < 0
	if !loopComplete {
		return current, false
	}

	if animator.Loops > 0 {
		animator.Loops--
	}

	if animator.Loops == 0 {
		if current >= frameCount {
			current = frameCount - 1
		} else if current < 0 {
			current = 0
		}
		return current, true
	}

	if animator.direction > 0 {
		current = 0
	} else {
		current = frameCount - 1
	}

	return current, false
}

func (m *AsepriteModule) getLibrary(entry *donburi.Entry) (*AnimationData, bool) {
	library, exists := m.animations[images.GetHandle(entry)]
	return library, exists
}

func (m *AsepriteModule) CreateSpriteEntity(ecs donburi.World, opts ...SpriteOption) *donburi.Entry {
	options := defaultSpriteOptions()
	for _, opt := range opts {
		opt(options)
	}

	entry := m.imageModule.CreateImageEntity(ecs,
		images.WithHandle(options.ImageHandle),
	)

	SetAnimationState(entry, options.State)
	SetAnimator(entry, &AnimatorInfo{
		Loops:     options.Loops,
		direction: 1,
	})
	SetClip(entry, options.Clip)

	return entry
}
