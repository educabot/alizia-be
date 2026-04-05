package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/educabot/team-ai-toolkit/boot"
	"github.com/educabot/team-ai-toolkit/dbconn"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/config"
	appweb "github.com/educabot/alizia-be/src/app/web"
)

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	server *boot.Server
}

func NewApp(cfg *config.Config) *App {
	db, err := dbconn.NewPostgresConnector(dbconn.PostgresConfig{
		URL: cfg.DatabaseURL,
	}).Connect()
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	repos := NewRepositories(db)
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
	err := sqlDB.Close()
	if err != nil {
		return
	}
}
