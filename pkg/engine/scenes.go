package engine

import (
	"context"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type (
	SceneExitCode    uint16
	SceneID          string
	SceneTransitions map[SceneExitCode]SceneID
)

type Scenes interface {
	AddScene(id SceneID, state SceneState, transitions SceneTransitions)
}

type SceneState struct {
	OnEnter  func(ctx context.Context, world donburi.World) error
	OnExit   func(ctx context.Context, world donburi.World) error
	OnUpdate func(ctx context.Context, world donburi.World) (SceneExitCode, error)
	OnRender func(ctx context.Context, world donburi.World, screen *ebiten.Image, viewMatrix ebiten.GeoM) error
}

type scenes struct {
	logger Logger

	scenes      map[SceneID]SceneState
	transitions map[SceneID]SceneTransitions

	currentScene SceneID
	nextScene    SceneID

	world donburi.World
}

const (
	SceneExitNone SceneExitCode = 0
	SceneIDNone   SceneID       = ""
)

func newScenes(initialScene SceneID, logger Logger) *scenes {
	return &scenes{
		logger:       logger,
		scenes:       make(map[SceneID]SceneState),
		transitions:  make(map[SceneID]SceneTransitions),
		currentScene: SceneIDNone,
		nextScene:    initialScene,
		world:        donburi.NewWorld(),
	}
}

func (s *scenes) AddScene(id SceneID, state SceneState, transitions SceneTransitions) {
	if _, exists := s.scenes[id]; exists {
		s.logger.Warn("Scene with ID '%s' is already registered. Overwriting.", id)
	}

	s.scenes[id] = state

	if _, exists := s.transitions[id]; exists {
		s.logger.Warn("Transitions for scene ID '%s' are already registered. Overwriting.", id)
	}

	s.transitions[id] = transitions
}

func (s *scenes) update(ctx context.Context) error {
	if s.currentScene == SceneIDNone && s.nextScene == SceneIDNone {
		return nil
	}

	if s.nextScene != SceneIDNone {
		return s.transitionToNext(ctx)
	}

	exitCode, err := s.updateCurrent(ctx)
	if err != nil {
		return err
	}

	if exitCode == SceneExitNone {
		return nil
	}

	transitions, ok := s.transitions[s.currentScene]
	if !ok {
		return fmt.Errorf("no transitions defined for current scene '%s'", s.currentScene)
	}

	nextScene, ok := transitions[exitCode]
	if !ok {
		return fmt.Errorf("no transition defined for exit code '%d' in scene '%s'", exitCode, s.currentScene)
	}

	s.nextScene = nextScene

	return nil
}

func (s *scenes) render(ctx context.Context, screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
	return s.renderCurrent(ctx, screen, viewMatrix)
}

func (s *scenes) transitionToNext(ctx context.Context) error {
	if err := s.exitCurrent(ctx); err != nil {
		return err
	}
	if err := s.enterNext(ctx); err != nil {
		return err
	}

	s.currentScene = s.nextScene
	s.nextScene = SceneIDNone

	return nil
}

func (s *scenes) exitCurrent(ctx context.Context) error {
	if s.currentScene == SceneIDNone {
		return nil
	}

	currentState, ok := s.scenes[s.currentScene]
	if !ok {
		return fmt.Errorf("scene with ID '%s' not found", s.currentScene)
	}

	if currentState.OnExit != nil {
		return currentState.OnExit(ctx, s.world)
	}

	return nil
}

func (s *scenes) enterNext(ctx context.Context) error {
	if s.nextScene == SceneIDNone {
		return nil
	}

	nextState, ok := s.scenes[s.nextScene]
	if !ok {
		return fmt.Errorf("scene with ID '%s' not found", s.nextScene)
	}

	if nextState.OnEnter != nil {
		return nextState.OnEnter(ctx, s.world)
	}

	return nil
}

func (s *scenes) updateCurrent(ctx context.Context) (SceneExitCode, error) {
	if s.currentScene == SceneIDNone {
		return SceneExitNone, nil
	}

	currentState, ok := s.scenes[s.currentScene]
	if !ok {
		return SceneExitNone, fmt.Errorf("scene with ID '%s' not found", s.currentScene)
	}

	if currentState.OnUpdate != nil {
		return currentState.OnUpdate(ctx, s.world)
	}

	return SceneExitNone, nil
}

func (s *scenes) renderCurrent(ctx context.Context, screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
	if s.currentScene == SceneIDNone {
		return nil
	}

	currentState, ok := s.scenes[s.currentScene]
	if !ok {
		return fmt.Errorf("scene with ID '%s' not found", s.currentScene)
	}

	if currentState.OnRender != nil {
		return currentState.OnRender(ctx, s.world, screen, viewMatrix)
	}

	return nil
}
