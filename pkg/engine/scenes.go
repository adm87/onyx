package engine

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type (
	SceneID       uint8
	SceneExitCode uint8
)

const (
	SceneIDNone       SceneID       = 0
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
	OnEnter func(donburi.World) error
	OnExit  func(donburi.World) error

	OnUpdate      []func(donburi.World, float64) (SceneExitCode, error)
	OnFixedUpdate []func(donburi.World, float64) error
	OnLateUpdate  []func(donburi.World, float64) error
	OnDraw        []func(donburi.World, *ebiten.Image) error
}

type sceneStateData struct {
	currentScene SceneID
	nextScene    SceneID
}

type Scenes interface {
	Register(id SceneID, def *SceneDefinition, transitions SceneTransitions)
	Start(id SceneID) error
	Update(deltaTime, fixedDeltaTime float64, fixedSteps int) error
	Draw(screen *ebiten.Image)
}

type scenes struct {
	world  donburi.World
	logger Logger

	sceneState *donburi.ComponentType[sceneStateData]
	stateQuery *donburi.Query

	definitions map[SceneID]*SceneDefinition
	transitions map[SceneID]SceneTransitions
}

func NewScenes(logger Logger) Scenes {
	world := donburi.NewWorld()

	sceneState := donburi.NewComponentType[sceneStateData]()
	world.Create(sceneState)

	return &scenes{
		world:       world,
		logger:      logger,
		sceneState:  sceneState,
		stateQuery:  donburi.NewQuery(filter.Contains(sceneState)),
		definitions: make(map[SceneID]*SceneDefinition),
		transitions: make(map[SceneID]SceneTransitions),
	}
}

func (s *scenes) Register(id SceneID, def *SceneDefinition, transitions SceneTransitions) {
	s.definitions[id] = def
	s.transitions[id] = transitions
}

func (s *scenes) Start(id SceneID) error {
	entry, ok := s.stateQuery.First(s.world)
	if !ok {
		return errors.New("scene state not found")
	}
	def, ok := s.definitions[id]
	if !ok {
		return fmt.Errorf("no definition found for scene %d", id)
	}
	if err := enterScene(def, s.world); err != nil {
		return err
	}
	s.sceneState.Get(entry).currentScene = id
	return nil
}

func (s *scenes) Update(deltaTime, fixedDeltaTime float64, fixedSteps int) error {
	entry, ok := s.stateQuery.First(s.world)
	if !ok {
		return errors.New("scene state not found")
	}

	state := s.sceneState.Get(entry)

	if !state.nextScene.IsNone() {
		current, hasCurrent := s.definitions[state.currentScene]
		next, hasNext := s.definitions[state.nextScene]
		if !hasNext {
			return fmt.Errorf("no definition found for scene %d", state.nextScene)
		}
		if hasCurrent {
			if err := exitScene(current, s.world); err != nil {
				return err
			}
		}
		if err := enterScene(next, s.world); err != nil {
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

	exitCode, err := runUpdates(def.OnUpdate, s.world, deltaTime)
	if err != nil {
		return err
	}

	if !exitCode.IsNone() {
		transitions, ok := s.transitions[state.currentScene]
		if !ok {
			s.logger.Warn("no transitions registered for scene %d", state.currentScene)
			return nil
		}
		next, ok := transitions[exitCode]
		if !ok {
			s.logger.Warn("no transition for exit code %d in scene %d", exitCode, state.currentScene)
			return nil
		}
		state.nextScene = next
		return nil
	}

	for range fixedSteps {
		if err := runFixed(def.OnFixedUpdate, s.world, fixedDeltaTime); err != nil {
			return err
		}
	}

	return runLate(def.OnLateUpdate, s.world, deltaTime)
}

func (s *scenes) Draw(screen *ebiten.Image) {
	entry, ok := s.stateQuery.First(s.world)
	if !ok {
		return
	}
	def, ok := s.definitions[s.sceneState.Get(entry).currentScene]
	if !ok {
		return
	}
	if err := runDraw(def.OnDraw, s.world, screen); err != nil {
		s.logger.Error("draw error: %v", err)
	}
}

func enterScene(def *SceneDefinition, world donburi.World) error {
	if def.OnEnter != nil {
		return def.OnEnter(world)
	}
	return nil
}

func exitScene(def *SceneDefinition, world donburi.World) error {
	if def.OnExit != nil {
		return def.OnExit(world)
	}
	return nil
}

func runUpdates(fns []func(donburi.World, float64) (SceneExitCode, error), world donburi.World, dt float64) (SceneExitCode, error) {
	for _, fn := range fns {
		code, err := fn(world, dt)
		if err != nil {
			return SceneExitCodeNone, err
		}
		if !code.IsNone() {
			return code, nil
		}
	}
	return SceneExitCodeNone, nil
}

func runFixed(fns []func(donburi.World, float64) error, world donburi.World, dt float64) error {
	for _, fn := range fns {
		if err := fn(world, dt); err != nil {
			return err
		}
	}
	return nil
}

func runLate(fns []func(donburi.World, float64) error, world donburi.World, dt float64) error {
	for _, fn := range fns {
		if err := fn(world, dt); err != nil {
			return err
		}
	}
	return nil
}

func runDraw(fns []func(donburi.World, *ebiten.Image) error, world donburi.World, screen *ebiten.Image) error {
	for _, fn := range fns {
		if err := fn(world, screen); err != nil {
			return err
		}
	}
	return nil
}
