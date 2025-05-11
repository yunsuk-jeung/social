package main

import (
	"time"

	"github.com/yunsuk-jeung/social/internal/db"
	"github.com/yunsuk-jeung/social/internal/env"
	"github.com/yunsuk-jeung/social/internal/mailer"
	"github.com/yunsuk-jeung/social/internal/store"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":3000"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:3000"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdelConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdelTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apikey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
	}

	// Logger
	logCfg := zap.NewDevelopmentConfig()
	logCfg.Encoding = "console"
	logCfg.DisableStacktrace = true
	logCfg.EncoderConfig.StacktraceKey = "stacktrace"
	logCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	prelogger, err := logCfg.Build()
	logger := prelogger.Sugar()

	if err != nil {
		return
	}

	// logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdelConns,
		cfg.db.maxIdelTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(db)

	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apikey, cfg.mail.fromEmail)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
