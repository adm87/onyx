package app

import (
	"context"
	"errors"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func Run(ctx context.Context, app Application, cfg *Config) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		initWindow(cfg)

		shl := &shell{
			ctx: NewContext(ctx, cfg,
				NewLogger(cfg.logOutput, cfg.logLevel),
				NewScreen(cfg.width, cfg.height),
			),
			app: app,
		}

		shl.ctx.lgr.Debug("application starting up...")
		shl.ctx.lgr.Debug("configuration: %+v", shl.ctx.cfg)

		if err := app.Startup(shl.ctx); err != nil {
			shl.ctx.lgr.Debug("application failed to start: %v", err)
			return err
		}

		runErr := ebiten.RunGame(shl)
		if runErr != nil {
			shl.ctx.lgr.Debug("application terminated with error: %v", runErr)
		}

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
		defer shutdownCancel()

		shl.ctx.lgr.Debug("application shutting down...")

		shutdownErr := app.Shutdown(shl.ctx.WithContext(shutdownCtx))
		if shutdownErr != nil {
			shl.ctx.lgr.Debug("application failed to shutdown: %v", shutdownErr)
		}

		return errors.Join(runErr, shutdownErr)
	}
}

func initWindow(cfg *Config) {
	title := cfg.name
	if cfg.version != "" {
		title += " " + cfg.version
	}
	ebiten.SetWindowTitle(title)
	ebiten.SetFullscreen(cfg.fullscreen)
	ebiten.SetWindowSize(cfg.width, cfg.height)
}

type shell struct {
	ctx *Context
	app Application
}

func (s *shell) Update() error {
	select {
	case <-s.ctx.Done():
		return s.ctx.ctx.Err()
	default:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return ebiten.Termination
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
			s.ctx.cfg.fullscreen = !s.ctx.cfg.fullscreen
			ebiten.SetFullscreen(s.ctx.cfg.fullscreen)
		}
		if err := s.app.Update(s.ctx); err != nil {
			s.ctx.Logger().Debug("application update error: %v", err)
			return err
		}
		return nil
	}
}

func (s *shell) Draw(screen *ebiten.Image) {
	select {
	case <-s.ctx.Done():
		return
	default:
		s.ctx.Screen().buffer.Fill(s.ctx.cfg.screenClearColor)

		if err := s.app.Draw(s.ctx); err != nil {
			s.ctx.Logger().Debug("application draw error: %v", err)
			return
		}

		screen.DrawImage(s.ctx.Screen().buffer, nil)
	}
}

func (s *shell) Layout(outsideWidth, outsideHeight int) (int, int) {
	return s.ctx.Screen().layout(outsideWidth, outsideHeight)
}
