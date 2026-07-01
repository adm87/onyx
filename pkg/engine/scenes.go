package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type (
	SceneExitCode    uint16
	SceneID          string
	SceneTransitions map[SceneExitCode]SceneID
)

type Scenes interface {
	AddScene(id SceneID, ctor SceneCtor, transitions SceneTransitions)
}

type Scene interface {
	Enter() error
	Exit() error

	Update(dt float64) (SceneExitCode, error)
	FixedUpdate(dt float64) error
	LateUpdate(dt float64) error

	Render(target *ebiten.Image) error
}

type SceneCtor func() Scene

type scenes struct {
	logger *logger

	currentScene SceneID
	nextScene    SceneID

	scenes      map[SceneID]SceneCtor
	transitions map[SceneID]SceneTransitions

	sceneInstance Scene
}

const (
	SceneExitNone SceneExitCode = 0
	SceneIDNone   SceneID       = ""
)

func newScenes(initialScene SceneID, logger *logger) *scenes {
	return &scenes{
		logger:       logger,
		currentScene: SceneIDNone,
		nextScene:    initialScene,
		scenes:       make(map[SceneID]SceneCtor),
		transitions:  make(map[SceneID]SceneTransitions),
	}
}

func (s *scenes) AddScene(id SceneID, ctor SceneCtor, transitions SceneTransitions) {
	if _, exists := s.scenes[id]; exists {
		s.logger.Warn("Scene with ID '%s' is already registered. Overwriting.", id)
	}

	s.scenes[id] = ctor

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

	exitCode, err := s.updateCurrent(deltaTime)
	if err != nil {
		return err
	}
	if exitCode != SceneExitNone {
		return s.setupNextTransition(exitCode)
	}
	if err := s.fixedUpdateCurrent(fixedDeltaTime, steps); err != nil {
		return err
	}

	return s.lateUpdateCurrent(deltaTime)
}

func (s *scenes) render(screen *ebiten.Image) error {
	return s.renderCurrent(screen)
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
	if s.sceneInstance != nil {
		return s.sceneInstance.Exit()
	}
	return nil
}

func (s *scenes) enterNext() error {
	if s.nextScene == SceneIDNone {
		return nil
	}

	nextScene, ok := s.scenes[s.nextScene]
	if !ok {
		return fmt.Errorf("scene with ID '%s' not found", s.nextScene)
	}

	s.sceneInstance = nextScene()
	if s.sceneInstance != nil {
		return s.sceneInstance.Enter()
	}
	return nil
}

func (s *scenes) updateCurrent(dt float64) (SceneExitCode, error) {
	if s.currentScene == SceneIDNone {
		return SceneExitNone, nil
	}
	if s.sceneInstance != nil {
		return s.sceneInstance.Update(dt)
	}
	return SceneExitNone, nil
}

func (s *scenes) fixedUpdateCurrent(dt float64, steps int) error {
	if s.currentScene == SceneIDNone {
		return nil
	}
	for range steps {
		if s.sceneInstance != nil {
			if err := s.sceneInstance.FixedUpdate(dt); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *scenes) lateUpdateCurrent(dt float64) error {
	if s.currentScene == SceneIDNone {
		return nil
	}
	if s.sceneInstance != nil {
		return s.sceneInstance.LateUpdate(dt)
	}
	return nil
}

func (s *scenes) renderCurrent(screen *ebiten.Image) error {
	if s.currentScene == SceneIDNone {
		return nil
	}
	if s.sceneInstance != nil {
		return s.sceneInstance.Render(screen)
	}
	return nil
}
