package engine

import (
	"context"
	"fmt"

	"github.com/adm87/onyx/pkg/engine/geom"
	"github.com/hajimehoshi/ebiten/v2"
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
	OnEnter       func(ctx context.Context, world World) error
	OnExit        func(ctx context.Context, world World) error
	OnUpdate      func(ctx context.Context, world World) (SceneExitCode, error)
	OnFixedUpdate func(ctx context.Context, world World) error
	OnLateUpdate  func(ctx context.Context, world World) error
	OnRender      func(ctx context.Context, world World, screen *ebiten.Image, viewMatrix ebiten.GeoM) error
}

type scenes struct {
	world  *world
	logger *logger

	currentScene SceneID
	nextScene    SceneID

	scenes      map[SceneID]SceneState
	transitions map[SceneID]SceneTransitions
}

const (
	SceneExitNone SceneExitCode = 0
	SceneIDNone   SceneID       = ""
)

func newScenes(initialScene SceneID, world *world, logger *logger) *scenes {
	return &scenes{
		world:        world,
		logger:       logger,
		currentScene: SceneIDNone,
		nextScene:    initialScene,
		scenes:       make(map[SceneID]SceneState),
		transitions:  make(map[SceneID]SceneTransitions),
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

func (s *scenes) update(ctx context.Context, steps int) error {
	if s.currentScene == SceneIDNone && s.nextScene == SceneIDNone {
		return nil
	}
	if s.nextScene != SceneIDNone {
		return s.transitionToNext(ctx)
	}

	currentState, ok := s.scenes[s.currentScene]
	if !ok {
		return fmt.Errorf("scene with ID '%s' not found", s.currentScene)
	}

	exitCode, err := s.updateCurrent(ctx, currentState)
	if err != nil {
		return err
	}
	if exitCode != SceneExitNone {
		return s.setupNextTransition(exitCode)
	}
	if err := s.fixedUpdateCurrent(ctx, currentState, steps); err != nil {
		return err
	}
	return s.lateUpdateCurrent(ctx, currentState)
}

func (s *scenes) render(ctx context.Context, region geom.AABB, screen *ebiten.Image, viewMatrix ebiten.GeoM) error {
	if err := s.world.render(screen, viewMatrix); err != nil {
		return err
	}
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

func (s *scenes) setupNextTransition(exitCode SceneExitCode) error {
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

func (s *scenes) updateCurrent(ctx context.Context, currentState SceneState) (SceneExitCode, error) {
	if s.currentScene == SceneIDNone {
		return SceneExitNone, nil
	}
	if currentState.OnUpdate != nil {
		return currentState.OnUpdate(ctx, s.world)
	}
	return SceneExitNone, nil
}

func (s *scenes) fixedUpdateCurrent(ctx context.Context, currentState SceneState, steps int) error {
	if s.currentScene == SceneIDNone {
		return nil
	}
	for range steps {
		if currentState.OnFixedUpdate != nil {
			if err := currentState.OnFixedUpdate(ctx, s.world); err != nil {
				return err
			}
		}
		if err := s.world.collision.checkCollision(s.world.ecs); err != nil {
			return err
		}
	}
	return nil
}

func (s *scenes) lateUpdateCurrent(ctx context.Context, currentState SceneState) error {
	if s.currentScene == SceneIDNone {
		return nil
	}
	if currentState.OnLateUpdate != nil {
		return currentState.OnLateUpdate(ctx, s.world)
	}
	return nil
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
