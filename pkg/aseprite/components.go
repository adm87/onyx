package aseprite

import (
	"time"

	"github.com/yohamta/donburi"
)

type State uint8

type AnimatorInfo struct {
	Loops     int
	direction int
	time      time.Duration
}

type SpriteOptions struct {
	State       State
	ImageHandle uint64
	Loops       int
	Frame       int
	Clip        string
}

type SpriteOption func(*SpriteOptions)

const (
	AnimationStateStopped State = iota
	AnimationStatePlaying
)

var (
	AnimationFrame = donburi.NewComponentType[int]()
	AnimationClip  = donburi.NewComponentType[string]()
	AnimationState = donburi.NewComponentType[State]()
	Animator       = donburi.NewComponentType[AnimatorInfo]()
)

func defaultSpriteOptions() *SpriteOptions {
	return &SpriteOptions{
		State:       AnimationStateStopped,
		ImageHandle: 0,
		Loops:       -1, // -1 for infinite loops
		Frame:       0,
		Clip:        "",
	}
}

func Playing() SpriteOption {
	return func(opts *SpriteOptions) {
		opts.State = AnimationStatePlaying
	}
}

func WithAnimationFrame(frame int) SpriteOption {
	return func(opts *SpriteOptions) {
		opts.Frame = frame
	}
}

func WithImageHandle(handle uint64) SpriteOption {
	return func(opts *SpriteOptions) {
		opts.ImageHandle = handle
	}
}

func WithLoops(loops int) SpriteOption {
	return func(opts *SpriteOptions) {
		opts.Loops = loops
	}
}

func WithClip(clip string) SpriteOption {
	return func(opts *SpriteOptions) {
		opts.Clip = clip
	}
}

func GetAnimator(entry *donburi.Entry) *AnimatorInfo {
	if !entry.HasComponent(Animator) {
		return nil
	}
	return Animator.Get(entry)
}

func SetAnimator(entry *donburi.Entry, info *AnimatorInfo) {
	donburi.Add(entry, Animator, info)
}

func GetAnimationState(entry *donburi.Entry) State {
	if !entry.HasComponent(AnimationState) {
		return AnimationStateStopped
	}
	return *AnimationState.Get(entry)
}

func SetAnimationState(entry *donburi.Entry, state State) {
	donburi.Add(entry, AnimationState, &state)
}

func GetLoops(entry *donburi.Entry) int {
	animator := GetAnimator(entry)
	if animator == nil {
		return 0
	}
	return animator.Loops
}

func SetLoops(entry *donburi.Entry, loops int) {
	animator := GetAnimator(entry)
	if animator == nil {
		return
	}
	animator.Loops = loops
	SetAnimator(entry, animator)
}

func GetAnimationFrame(entry *donburi.Entry) int {
	if !entry.HasComponent(AnimationFrame) {
		return 0
	}
	return *AnimationFrame.Get(entry)
}

func SetAnimationFrame(entry *donburi.Entry, frame int) {
	donburi.Add(entry, AnimationFrame, &frame)
}

func GetClip(entry *donburi.Entry) string {
	if !entry.HasComponent(AnimationClip) {
		return ""
	}
	return *AnimationClip.Get(entry)
}

func SetClip(entry *donburi.Entry, clip string) {
	curr := GetClip(entry)
	if curr == clip {
		return
	}

	animator := GetAnimator(entry)
	if animator != nil {
		animator.time = 0
		SetAnimator(entry, animator)
		SetAnimationState(entry, AnimationStatePlaying)
	}

	SetAnimationFrame(entry, 0)
	donburi.Add(entry, AnimationClip, &clip)
}

func IsIdle(entry *donburi.Entry) bool {
	state := GetAnimationState(entry)
	return state == AnimationStateStopped
}

func IsPlaying(entry *donburi.Entry) bool {
	state := GetAnimationState(entry)
	return state == AnimationStatePlaying
}
