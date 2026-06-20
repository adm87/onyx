package aseprite

// func (m *AsepritePlugin) UpdateAnimation(entry *donburi.Entry, dt time.Duration) {
// 	animator := GetAnimator(entry)
// 	if animator == nil {
// 		return
// 	}

// 	if animator.State != AnimationStatePlaying {
// 		return
// 	}

// 	library, exists := m.getLibrary(entry)
// 	if !exists {
// 		return
// 	}

// 	tag, exists := library.Meta.Clips[animator.Clip]
// 	if !exists {
// 		return
// 	}

// 	elapsed := animator.time + dt

// 	frame := animator.Frame
// 	frameIndex := tag.From + frame

// 	duration := time.Duration(library.Frames[frameIndex].Duration) * time.Millisecond
// 	if duration <= 0 || elapsed < duration {
// 		animator.time = elapsed
// 		return
// 	}

// 	nextFrame := frame
// 	frameCount := tag.To - tag.From + 1

// 	var completed bool
// 	for elapsed >= duration {
// 		elapsed -= duration

// 		nextFrame, completed = m.getNextFrame(nextFrame, animator, frameCount)
// 		if completed {
// 			break
// 		}

// 		frameIndex = tag.From + nextFrame
// 		duration = time.Duration(library.Frames[frameIndex].Duration) * time.Millisecond
// 	}

// 	if nextFrame != frame {
// 		images.SetFrame(entry, frameIndex)
// 		animator.Frame = nextFrame
// 	}
// 	if animator.time != elapsed {
// 		animator.time = elapsed
// 	}
// 	if completed {
// 		animator.State = AnimationStateStopped
// 	}
// }

// func (m *AsepritePlugin) getNextFrame(current int, animator *AnimatorModel, frameCount int) (int, bool) {
// 	current += animator.direction

// 	loopComplete := current >= frameCount || current < 0
// 	if !loopComplete {
// 		return current, false
// 	}

// 	if animator.Loops > 0 {
// 		animator.Loops--
// 	}

// 	if animator.Loops == 0 {
// 		if current >= frameCount {
// 			current = frameCount - 1
// 		} else if current < 0 {
// 			current = 0
// 		}
// 		return current, true
// 	}

// 	if animator.direction > 0 {
// 		current = 0
// 	} else {
// 		current = frameCount - 1
// 	}

// 	return current, false
// }

// func (m *AsepritePlugin) getLibrary(entry *donburi.Entry) (*AnimationData, bool) {
// 	library, exists := m.animations[images.GetHandle(entry)]
// 	return library, exists
// }

// func (m *AsepritePlugin) CreateSpriteEntity(ecs donburi.World, opts ...SpriteOption) *donburi.Entry {
// 	options := defaultSpriteOptions()
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	entry := m.imagesPlugin.CreateImageEntity(ecs,
// 		images.WithHandle(options.ImageHandle),
// 	)

// 	SetAnimationState(entry, options.State)
// 	SetAnimator(entry, &AnimatorModel{
// 		Loops:     options.Loops,
// 		direction: 1,
// 	})
// 	SetClip(entry, options.Clip)

// 	return entry
// }
