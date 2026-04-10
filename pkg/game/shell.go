package game

import (
	"errors"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameFunc func(ctx Context) error

func gameNoop(ctx Context) error {
	return nil
}

type Shell struct {
	opts *Options
	ctx  *ctxImpl
}

func NewShell(opts ...Option) *Shell {
	o := defaultOpts()
	applyOpts(o, opts...)
	return &Shell{
		opts: o,
		ctx: &ctxImpl{
			ctx: o.ctx,
			mdl: o.mdl,
			lgr: o.lgr,
			srn: &screenImpl{
				img: ebiten.NewImage(o.width, o.height),
				opt: &ebiten.DrawImageOptions{
					Filter: o.screenFilter,
				},
			},
			tm: &timeImpl{
				fixedDelta: time.Duration(time.Second / time.Duration(o.fps)),
			},
		},
	}
}

func (s *Shell) Context() Context {
	return s.ctx
}

func (s *Shell) Run() error {
	select {
	case <-s.opts.ctx.Done():
		return s.opts.ctx.Err()
	default:
		s.initWindow()

		if err := s.opts.OnStart(s.ctx); err != nil {
			return fmt.Errorf("failed to start game: %w", err)
		}

		return ebiten.RunGameWithOptions(s, &s.opts.runOpts)
	}
}

func (s *Shell) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.ctx.srn.layout(outsideWidth, outsideHeight)
}

func (s *Shell) Update() error {
	select {
	case <-s.opts.ctx.Done():
		return s.opts.ctx.Err()
	default:
		s.ctx.tm.tick()

		if err := s.updateGame(); err != nil {
			if errors.Is(err, ebiten.Termination) {
				return err
			}
			fmt.Printf("error during update: %v\n", err)
			return err
		}
		return nil
	}
}

func (s *Shell) Draw(screen *ebiten.Image) {
	select {
	case <-s.opts.ctx.Done():
		return
	default:
		s.ctx.srn.img.Fill(s.opts.clearColor)

		if err := s.opts.OnDraw(s.ctx); err != nil {
			fmt.Printf("error during draw: %v\n", err)
			return
		}

		screen.DrawImage(s.ctx.srn.img, nil)
	}
}

func (s *Shell) initWindow() {
	ebiten.SetWindowTitle(s.opts.title)
	ebiten.SetWindowSize(s.opts.width, s.opts.height)
	ebiten.SetFullscreen(s.opts.fullscreen)
}

func (s *Shell) updateGame() error {
	if err := s.opts.OnUpdate(s.ctx); err != nil {
		if errors.Is(err, ebiten.Termination) {
			return err
		}
		return fmt.Errorf("error during update: %w", err)
	}
	for i := 0; i < s.ctx.tm.steps; i++ {
		if err := s.opts.OnFixedUpdate(s.ctx); err != nil {
			return fmt.Errorf("error during fixed update: %w", err)
		}
	}
	if err := s.opts.OnLateUpdate(s.ctx); err != nil {
		return fmt.Errorf("error during late update: %w", err)
	}
	return nil
}
