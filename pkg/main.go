package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"github.com/fatih/color"

	pkg "git.solsynth.dev/hypernet/passport/pkg/internal"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"

	"git.solsynth.dev/hypernet/passport/pkg/internal/grpc"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web"
	"github.com/robfig/cron/v3"

	"git.solsynth.dev/hypernet/passport/pkg/internal/cache"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	// Booting screen
	fmt.Println(color.YellowString(" ____                                _\n|  _ \\ __ _ ___ ___ _ __   ___  _ __| |_\n| |_) / _` / __/ __| '_ \\ / _ \\| '__| __|\n|  __/ (_| \\__ \\__ \\ |_) | (_) | |  | |_\n|_|   \\__,_|___/___/ .__/ \\___/|_|   \\__|\n                   |_|"))
	fmt.Printf("%s v%s\n", color.New(color.FgHiYellow).Add(color.Bold).Sprintf("Hypernet.Passport"), pkg.AppVersion)
	fmt.Printf("The user identity service in Hypernet\n")
	color.HiBlack("=====================================================\n")

	// Configure settings
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")

	// Load settings
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading settings.")
	}

	// Connect to nexus
	if err := gap.InitializeToNexus(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connecting to nexus...")
	}

	// Load keypair
	if reader, err := sec.NewInternalTokenReader(viper.GetString("security.internal_public_key")); err != nil {
		log.Error().Err(err).Msg("An error occurred when reading internal public key for jwt. Authentication related features will be disabled.")
	} else {
		web.IReader = reader
		log.Info().Msg("Internal jwt public key loaded.")
	}
	if reader, err := sec.NewJwtReader(viper.GetString("security.public_key")); err != nil {
		log.Error().Err(err).Msg("An error occurred when reading public key for jwt. Signing token may not work.")
	} else {
		services.EReader = reader
		log.Info().Msg("Jwt public key loaded.")
	}
	if writer, err := sec.NewJwtWriter(viper.GetString("security.private_key")); err != nil {
		log.Error().Err(err).Msg("An error occurred when reading private key for jwt. Signing token may not work.")
	} else {
		services.EWriter = writer
		log.Info().Msg("Jwt private key loaded.")
	}

	// Load localization
	if err := gap.LoadLocalization(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when loading localization.")
	}

	// Connect to database
	if err := database.NewGorm(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connect to database.")
	} else if err := database.RunMigration(database.C); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when running database auto migration.")
	}
	if err := database.NewGeoDB(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connect to geoip database.")
	}

	// Initialize cache
	if err := cache.NewStore(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when initializing cache.")
	}

	// App
	go web.NewServer().Listen()

	// Grpc App
	go grpc.NewServer().Listen()

	// Configure timed tasks
	quartz := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(&log.Logger)))
	quartz.AddFunc("@every 60m", services.DoAutoSignoff)
	quartz.AddFunc("@every 60m", services.DoAutoDatabaseCleanup)
	quartz.AddFunc("@midnight", services.RecycleUnConfirmAccount)
	quartz.AddFunc("@every 60s", services.SaveEventChanges)
	quartz.Start()

	// Messages
	log.Info().Msgf("Passport v%s is started...", pkg.AppVersion)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msgf("Passport v%s is quitting...", pkg.AppVersion)

	quartz.Stop()
}
