package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/contrib/bridges/otelslog"

	"github.com/trollLemon/DiscordBot/internal/application"
	"github.com/trollLemon/DiscordBot/internal/commands"
	"github.com/trollLemon/DiscordBot/internal/common"
	"github.com/trollLemon/DiscordBot/internal/telemetry"
)

type Options struct {
	RegisterCommands bool
}

func parseCommandLineArgs() *Options {

	shouldRegisterCommands := flag.Bool("register-commands", true, "register bot commands to guild")
	flag.Parse()

	return &Options{
		*shouldRegisterCommands,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		slog.Warn("Invalid duration in env var, using default", "env", key, "value", raw, "error", err)
		return defaultValue
	}

	return value
}

func main() {
	options := parseCommandLineArgs()

	exporterEndpoint := getEnvOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318")
	shutdownTimeout := getEnvDurationOrDefault("OTEL_SHUTDOWN_TIMEOUT", 5*time.Second)

	shutdownTracer, err := telemetry.InitTracer(exporterEndpoint)
	if err != nil {
		slog.Error("Failed to initialize tracing", "error", err)
		os.Exit(1)
	}

	shutdownLogger, err := telemetry.InitLogger(exporterEndpoint)
	if err != nil {
		slog.Error("Failed to initialize OTel logger", "error", err)
		os.Exit(1)
	}

	slog.SetDefault(otelslog.NewLogger("discord-bot"))

	slog.Info("Reading config from ENV")

	conf, err := common.GetBotConfig()
	if err != nil {
		slog.Error("invalid bot configuration", "error", err)
		os.Exit(1)
	}

	session, err := discordgo.New("Bot " + conf.BotToken)
	if err != nil {
		slog.Error("error creating Discord session", "error", err)
		os.Exit(1)
	}

	slog.Info("Created a Discord Session")

	if err := session.Open(); err != nil {
		slog.Error("error opening connection", "error", err)
		os.Exit(1)
	}

	slog.Info("Connected to Discord")

	if options.RegisterCommands {
		go func() {
			slog.Info("Registering commands...")
			commands.RegisterCommands(session)
		}()
	}

	slog.Info("Initializing application")
	app := application.InitializeApplication(conf, context.Background())

	commands.AddCommandHandlers(session, app)

	slog.Info("Bot is online.")
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	slog.Info("Press Ctrl+C to stop the bot")
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := shutdownTracer(shutdownCtx); err != nil {
		slog.Error("Failed to flush tracing data during shutdown", "error", err)
	}
	if err := shutdownLogger(shutdownCtx); err != nil {
		slog.Error("Failed to flush logging data during shutdown", "error", err)
	}
}
