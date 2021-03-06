package main

import (
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"majipoor/cmd/majipoor/gen"
	"majipoor/cmd/majipoor/mysql"
	"majipoor/lib/helpers"
	"os"
	"strings"
	"time"
)

var rootCmd = cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		if viper.GetBool("log.debug") {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		file := viper.GetString("log.log-file")
		if file == "" {
			if isatty.IsTerminal(os.Stderr.Fd()) {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
			} else {
				log.Logger = log.Output(os.Stderr)
			}
		} else {
			w, err := os.Open(file)
			if err != nil {
				log.Fatal().Err(err).Msgf("Could not open log file %s", file)
			}
			log.Debug().Str("log-file", file).Msg("Logging to file")
			log.Logger = log.Output(w)
		}

		if viper.GetBool("log.log-line") {
			log.Logger = log.With().Caller().Logger()
		}

		if viper.GetBool("log.log-error-stacktrace") {
			log.Debug().Msg("Logging error stacktraces")
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		}

		if viper.GetBool("log.log-goroutines") {
			helpers.StartBackgroundGoroutinePrinter()
		}
	},
}

func viperBindNestedPFlags(namespace string, cmd *cobra.Command, flags []string) error {
	for _, flag := range flags {
		f := cmd.PersistentFlags().Lookup(flag)
		viperFlagName := namespace + "." + flag
		if strings.HasPrefix(flag, namespace+"-") {
			viperFlagName = namespace + "." + flag[len(namespace)+1:]
		}
		if err := viper.BindPFlag(viperFlagName, f); err != nil {
			return errors.Wrapf(err, "Could not bind flag %s to viper flag %s", flag, viperFlagName)
		}
	}

	return nil
}

func main() {
	viper.SetConfigName("majipoor")
	viper.AddConfigPath("$HOME/.config")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := viper.ReadInConfig(); err != nil {
		// never show this error for now
		log.Trace().Err(err).Msg("Failed to read config")
	}

	rootCmd.PersistentFlags().Bool("log-goroutines", false, "Periodically log goroutine count")
	rootCmd.PersistentFlags().Bool("log-debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().Bool("log-error-stacktrace", false, "Enable stacktrace logging on errors")
	rootCmd.PersistentFlags().Bool("log-line", true, "Enable logging of file ane line number")
	rootCmd.PersistentFlags().String("log-file", "", "Enable logging to file")
	if err := viperBindNestedPFlags("log", &rootCmd,
		[]string{"log-debug", "log-error-stacktrace", "log-line", "log-file", "log-goroutines"}); err != nil {
		log.Fatal().Err(err).Msg("Could not bind persistent flags")
	}

	rootCmd.PersistentFlags().String("mysql-host", "localhost", "Mysql hostname")
	rootCmd.PersistentFlags().String("mysql-username", "mysql", "Mysql username")
	rootCmd.PersistentFlags().String("mysql-password", "", "Mysql password")
	rootCmd.PersistentFlags().Int("mysql-port", 3306, "Mysql port")
	rootCmd.PersistentFlags().String("mysql-db", "", "Mysql database")
	rootCmd.PersistentFlags().StringArray("mysql-limit-tables", []string{}, "Restrict syncing to these tables (if not empty)")
	rootCmd.PersistentFlags().StringArray("mysql-skip-tables", []string{}, "Skip these tables when syncing")
	rootCmd.PersistentFlags().String("mysql-root-username", "root", "Mysql root username")
	rootCmd.PersistentFlags().String("mysql-root-password", "master", "Mysql root password")
	if err := viperBindNestedPFlags("mysql", &rootCmd,
		[]string{"mysql-host", "mysql-username", "mysql-password", "mysql-port", "mysql-db",
			"mysql-limit-tables", "mysql-skip-tables",
			"mysql-root-username", "mysql-root-password"}); err != nil {
		log.Fatal().Err(err).Msg("Could not bind persistent flags")
	}

	rootCmd.PersistentFlags().String("postgresql-host", "localhost", "PG hostname")
	rootCmd.PersistentFlags().String("postgresql-username", "postgres", "PG username")
	rootCmd.PersistentFlags().String("postgresql-password", "master", "PG password")
	rootCmd.PersistentFlags().Int("postgresql-port", 5432, "PG port")
	rootCmd.PersistentFlags().String("postgresql-db", "postgres", "PG database")
	rootCmd.PersistentFlags().String("postgresql-schema", "majipoor", "PG destination schema")
	if err := viperBindNestedPFlags("postgresql", &rootCmd,
		[]string{"postgresql-host", "postgresql-username", "postgresql-password", "postgresql-port", "postgresql-db", "postgresql-schema"}); err != nil {
		log.Fatal().Err(err).Msg("Could not bind persistent flags")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix("MAJIPOOR")
	viper.AutomaticEnv()

	rootCmd.AddCommand(mysql.MysqlCmd)
	rootCmd.AddCommand(gen.GenCmd)

	helpers.StartSIGPROFStacktraceDumper("")

	_ = rootCmd.Execute()
}
