package engine

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type (
	SceneID       string
	SceneExitCode uint8
)

const (
	SceneIDNone       SceneID       = ""
	SceneExitCodeNone SceneExitCode = 0
)

func (id SceneID) IsNone() bool {
	return id == SceneIDNone
}

func (code SceneExitCode) IsNone() bool {
	return code == SceneExitCodeNone
}

type SceneTransitions map[SceneExitCode]SceneID

type SceneDefinition struct {
	SceneID SceneID

	OnEnter       func(Scene) error
	OnExit        func(Scene) error
	OnUpdate      func(Scene, float64) (SceneExitCode, error)
	OnFixedUpdate func(Scene, float64) error
	OnLateUpdate  func(Scene, float64) error
	OnDraw        func(Scene, *ebiten.Image) error
}

type sceneStateData struct {
	currentScene SceneID
	nextScene    SceneID
}

type Scenes interface {
	RegisterScenes(defs ...*SceneDefinition) error
	RegisterTransitions(sceneID SceneID, transitions SceneTransitions) error
	Start(id SceneID) error
	Update(deltaTime, fixedDeltaTime float64, fixedSteps int) error
	Draw(screen *ebiten.Image)
}

type scenes struct {
	scene  Scene
	logger Logger

	sceneState *donburi.ComponentType[sceneStateData]
	stateQuery *donburi.Query

	definitions map[SceneID]*SceneDefinition
	transitions map[SceneID]SceneTransitions
}

func NewScenes(logger Logger) Scenes {
	scene := newScene()

	sceneState := donburi.NewComponentType[sceneStateData]()
	scene.World().Create(sceneState)

	return &scenes{
		scene:       scene,
		logger:      logger,
		sceneState:  sceneState,
		stateQuery:  donburi.NewQuery(filter.Contains(sceneState)),
		definitions: make(map[SceneID]*SceneDefinition),
		transitions: make(map[SceneID]SceneTransitions),
	}
}

func (s *scenes) RegisterScenes(defs ...*SceneDefinition) error {
	for _, def := range defs {
		if def == nil {
			return errors.New("scene definition cannot be nil")
		}
		if def.SceneID.IsNone() {
			return errors.New("scene definition must have a valid SceneID")
		}
		if _, exists := s.definitions[def.SceneID]; exists {
			return fmt.Errorf("scene definition with ID %s already exists", def.SceneID)
		}
		s.definitions[def.SceneID] = def
	}
	return nil
}

func (s *scenes) RegisterTransitions(sceneID SceneID, transitions SceneTransitions) error {
	if _, exists := s.definitions[sceneID]; !exists {
		return fmt.Errorf("cannot register transitions for undefined scene ID %s", sceneID)
	}
	s.transitions[sceneID] = transitions
	return nil
}

func (s *scenes) Start(id SceneID) error {
	entry, ok := s.stateQuery.First(s.scene.World())
	if !ok {
		return errors.New("scene state not found")
	}
	def, ok := s.definitions[id]
	if !ok {
		return fmt.Errorf("no definition found for scene %s", id)
	}
	if err := enterScene(def, s.scene); err != nil {
		return err
	}
	s.sceneState.Get(entry).currentScene = id
	return nil
}

func (s *scenes) Update(deltaTime, fixedDeltaTime float64, fixedSteps int) error {
	world := s.scene.World()

	entry, ok := s.stateQuery.First(world)
	if !ok {
		return errors.New("scene state not found")
	}

	state := s.sceneState.Get(entry)

	if !state.nextScene.IsNone() {
		current, hasCurrent := s.definitions[state.currentScene]
		next, hasNext := s.definitions[state.nextScene]

		if !hasNext {
			return fmt.Errorf("no definition found for scene %s", state.nextScene)
		}
		if hasCurrent {
			if err := exitScene(current, s.scene); err != nil {
				return err
			}
		}
		if err := enterScene(next, s.scene); err != nil {
			return err
		}

		state.currentScene = state.nextScene
		state.nextScene = SceneIDNone
		return nil
	}

	def, ok := s.definitions[state.currentScene]
	if !ok {
		return nil
	}

	if update := def.OnUpdate; update != nil {
		exitCode, err := update(s.scene, deltaTime)
		if err != nil {
			return err
		}

		if !exitCode.IsNone() {
			transitions, ok := s.transitions[state.currentScene]
			if !ok {
				s.logger.Warn("no transitions registered for scene %s", state.currentScene)
				return nil
			}
			next, ok := transitions[exitCode]
			if !ok {
				s.logger.Warn("no transition for exit code %d in scene %s", exitCode, state.currentScene)
				return nil
			}
			state.nextScene = next
			return nil
		}
	}

	if fixedUpdate := def.OnFixedUpdate; fixedUpdate != nil {
		for range fixedSteps {
			if err := fixedUpdate(s.scene, fixedDeltaTime); err != nil {
				return err
			}
		}
	}

	if lateUpdate := def.OnLateUpdate; lateUpdate != nil {
		if err := lateUpdate(s.scene, deltaTime); err != nil {
			return err
		}
	}

	return nil
}

func (s *scenes) Draw(screen *ebiten.Image) {
	entry, ok := s.stateQuery.First(s.scene.World())
	if !ok {
		return
	}
	def, ok := s.definitions[s.sceneState.Get(entry).currentScene]
	if !ok {
		return
	}
	if draw := def.OnDraw; draw != nil {
		if err := draw(s.scene, screen); err != nil {
			s.logger.Error("error drawing scene %s: %v", def.SceneID, err)
		}
	}
}

func enterScene(def *SceneDefinition, world Scene) error {
	if def.OnEnter != nil {
		return def.OnEnter(world)
	}
	return nil
}

func exitScene(def *SceneDefinition, world Scene) error {
	if def.OnExit != nil {
		return def.OnExit(world)
	}
	return nil
}
