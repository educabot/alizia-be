package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/educabot/alizia-be/config"
	appweb "github.com/educabot/alizia-be/src/app/web"
	"github.com/educabot/team-ai-toolkit/boot"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	server *boot.Server
}

func NewApp(cfg *config.Config, db *gorm.DB) *App {
	repos := NewRepositories(cfg, db)
	usecases := NewUseCases(repos)
	container := NewHandlers(usecases, cfg)

	engine := boot.NewRouter(cfg.Env, cfg.AllowedOrigins)

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	appweb.ConfigureMappings(engine, container, cfg)

	server := boot.NewServer(cfg.Port, engine)

	return &App{cfg: cfg, db: db, server: server}
}

func (a *App) Run() {
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		a.server.Shutdown()
	}()

	a.server.Run()
}

func (a *App) Close() {
	sqlDB, _ := a.db.DB()
	sqlDB.Close()
}
