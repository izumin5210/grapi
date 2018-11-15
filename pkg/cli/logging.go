package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggingMode int

const (
	loggingNop loggingMode = iota
	loggingVerbose
	loggingDebug
)

var (
	logging = loggingNop

	// DebugLogConfig is used to generate a *zap.Logger for debug mode.
	DebugLogConfig = func() zap.Config {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.DisableStacktrace = true
		return cfg
	}()
	// VerboseLogConfig is used to generate a *zap.Logger for verbose mode.
	VerboseLogConfig = func() zap.Config {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Local().Format("2006-01-02 15:04:05 MST"))
		}
		return cfg
	}()

	closers []func()
)

// AddLoggingFlags sets "--debug" and "--verbose" flags to the given *cobra.Command instance.
func AddLoggingFlags(cmd *cobra.Command) {
	var (
		debugEnabled, verboseEnabled bool
	)

	cmd.PersistentFlags().BoolVar(
		&debugEnabled,
		"debug",
		false,
		fmt.Sprintf("Debug level output"),
	)
	cmd.PersistentFlags().BoolVarP(
		&verboseEnabled,
		"verbose",
		"v",
		false,
		fmt.Sprintf("Verbose level output"),
	)

	cobra.OnInitialize(func() {
		switch {
		case debugEnabled:
			Debug()
		case verboseEnabled:
			Verbose()
		}
	})
}

// Debug sets a debug logger in global.
func Debug() {
	logging = loggingDebug
	replaceLogger(DebugLogConfig)
}

// Verbose sets a verbose logger in global.
func Verbose() {
	logging = loggingVerbose
	replaceLogger(VerboseLogConfig)
}

// IsDebug returns true if a debug logger is used.
func IsDebug() bool { return logging == loggingDebug }

// IsVerbose returns true if a verbose logger is used.
func IsVerbose() bool { return logging == loggingVerbose }

func replaceLogger(cfg zap.Config) {
	l, err := cfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize a debug logger: %v\n", err)
	}

	closers = append(closers, func() { l.Sync() })
	closers = append(closers, zap.ReplaceGlobals(l))
}

// Close closes cli utilities.
func Close() {
	for _, f := range closers {
		f()
	}
}
