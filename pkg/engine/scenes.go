package engine

import (
	"fmt"

	"github.com/adm87/onyx/pkg/engine/geom"
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
	OnEnter       func(ecs donburi.World) error
	OnExit        func(ecs donburi.World) error
	OnUpdate      func(ecs donburi.World, dt float64) (SceneExitCode, error)
	OnFixedUpdate func(ecs donburi.World, dt float64) error
	OnLateUpdate  func(ecs donburi.World, dt float64) error
	OnRender      func(ecs donburi.World, screen *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) error
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

func (s *scenes) update(steps int, deltaTime float64, fixedDeltaTime float64) error {
	if s.currentScene == SceneIDNone && s.nextScene == SceneIDNone {
		return nil
	}
	if s.nextScene != SceneIDNone {
		return s.transitionToNext()
	}

	currentState, ok := s.scenes[s.currentScene]
	if !ok {
		return fmt.Errorf("scene with ID '%s' not found", s.currentScene)
	}

	exitCode, err := s.updateCurrent(currentState, deltaTime)
	if err != nil {
		return err
	}
	if exitCode != SceneExitNone {
		return s.setupNextTransition(exitCode)
	}
	if err := s.fixedUpdateCurrent(currentState, fixedDeltaTime, steps); err != nil {
		return err
	}
	return s.lateUpdateCurrent(currentState, deltaTime)
}

func (s *scenes) render(screen *ebiten.Image, viewPort geom.AABB, viewMatrix ebiten.GeoM) error {
	return s.renderCurrent(screen, viewPort, viewMatrix)
}

func (s *scenes) transitionToNext() error {
	if err := s.exitCurrent(); err != nil {
		return err
	}
	if err := s.enterNext(); err != nil {
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

func (s *scenes) exitCurrent() error {
	if s.currentScene == SceneIDNone {
		return nil
	}

	currentState, ok := s.scenes[s.currentScene]
	if !ok {
		return fmt.Errorf("scene with ID '%s' not found", s.currentScene)
	}
	if currentState.OnExit != nil {
		return currentState.OnExit(s.world.ecs)
	}
	return nil
}

func (s *scenes) enterNext() error {
	if s.nextScene == SceneIDNone {
		return nil
	}

	nextState, ok := s.scenes[s.nextScene]
	if !ok {
		return fmt.Errorf("scene with ID '%s' not found", s.nextScene)
	}
	if nextState.OnEnter != nil {
		return nextState.OnEnter(s.world.ecs)
	}
	return nil
}

func (s *scenes) updateCurrent(currentState SceneState, dt float64) (SceneExitCode, error) {
	if s.currentScene == SceneIDNone {
		return SceneExitNone, nil
	}
	if currentState.OnUpdate != nil {
		return currentState.OnUpdate(s.world.ecs, dt)
	}
	return SceneExitNone, nil
}

func (s *scenes) fixedUpdateCurrent(currentState SceneState, dt float64, steps int) error {
	if s.currentScene == SceneIDNone {
		return nil
	}
	for range steps {
		if currentState.OnFixedUpdate != nil {
			if err := currentState.OnFixedUpdate(s.world.ecs, dt); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *scenes) lateUpdateCurrent(currentState SceneState, dt float64) error {
	if s.currentScene == SceneIDNone {
		return nil
	}
	if currentState.OnLateUpdate != nil {
		return currentState.OnLateUpdate(s.world.ecs, dt)
	}
	return nil
}

func (s *scenes) renderCurrent(screen *ebiten.Image, viewport geom.AABB, viewMatrix ebiten.GeoM) error {
	if s.currentScene == SceneIDNone {
		return nil
	}

	currentState, ok := s.scenes[s.currentScene]
	if !ok {
		return fmt.Errorf("scene with ID '%s' not found", s.currentScene)
	}
	if currentState.OnRender != nil {
		return currentState.OnRender(s.world.ecs, screen, viewport, viewMatrix)
	}
	return nil
}
