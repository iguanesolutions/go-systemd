package sysdjournaldslog

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	sysdjournald "github.com/iguanesolutions/go-systemd/v5/journald"
)

const (
	// https://wiki.archlinux.org/title/Systemd/Journal#Priority_level
	LevelDebug       = slog.LevelDebug
	LevelDebugStr    = "DEBUG"
	LevelInfo        = slog.LevelInfo
	LevelInfoStr     = "INFO"
	LevelNotice      = LevelInfo + 2
	LevelNoticeStr   = "NOTICE"
	LevelWarning     = slog.LevelWarn
	LevelWarningStr  = "WARNING"
	LevelError       = slog.LevelError
	LevelErrorStr    = "ERROR"
	LevelCritical    = LevelError + 2
	LevelCriticalStr = "CRITICAL"
	LevelAlert       = LevelCritical + 2
	LevelAlertStr    = "ALERT"
	// LevelEmergency should not be used by applications
	LevelEmergency    = LevelAlert + 2
	LevelEmergencyStr = "EMERGENCY"
)

const (
	prefixDebugStr     = sysdjournald.DebugPrefix + slog.LevelKey
	prefixInfoStr      = sysdjournald.InfoPrefix + slog.LevelKey
	prefixNoticeStr    = sysdjournald.NoticePrefix + slog.LevelKey
	prefixWarningStr   = sysdjournald.WarningPrefix + slog.LevelKey
	prefixErrorStr     = sysdjournald.ErrPrefix + slog.LevelKey
	prefixCriticalStr  = sysdjournald.CritPrefix + slog.LevelKey
	prefixAlertStr     = sysdjournald.AlertPrefix + slog.LevelKey
	prefixEmergencyStr = sysdjournald.EmergPrefix + slog.LevelKey
)

// GetAvailableLogLevels returns a list of available log levels that can be used by GetLogLevel()
func GetAvailableLogLevels() []string {
	return []string{
		LevelDebugStr,
		LevelInfoStr,
		LevelNoticeStr,
		LevelWarningStr,
		LevelErrorStr,
		LevelCriticalStr,
		LevelAlertStr,
		LevelEmergencyStr,
	}
}

// GetLogLevel returns a log level based on the given string. If the string is not recognized, it will return LevelInfo.
func GetLogLevel(raw string) slog.Leveler {
	switch strings.ToUpper(raw) {
	case LevelDebugStr:
		return LevelDebug
	case LevelInfoStr:
		return LevelInfo
	case LevelNoticeStr:
		return LevelNotice
	case LevelWarningStr:
		return LevelWarning
	case LevelErrorStr:
		return LevelError
	case LevelCriticalStr:
		return LevelCritical
	case LevelAlertStr:
		return LevelAlert
	case LevelEmergencyStr:
		return LevelEmergency
	default:
		return LevelInfo
	}
}

// Options represents the options for the journald slog handler.
type Options struct {
	// AddSource causes the handler to compute the source code position
	// of the log statement and add a SourceKey attribute to the output.
	AddSource bool
	// SourceFormat specifies a function that formats the source information.
	SourceFormat func(*slog.Source) string
	// Level reports the minimum record level that will be logged.
	// The handler discards records with lower levels.
	// If Level is nil, the handler assumes LevelInfo.
	// The handler calls Level.Level for each record processed;
	// to adjust the minimum level dynamically, use a LevelVar.
	Level slog.Leveler
}

// NewHandler returns a new slog handler that writes logs in a journald compatible/enhanced format.
func NewHandler(opts Options) slog.Handler {
	if opts.SourceFormat == nil {
		opts.SourceFormat = func(src *slog.Source) string {
			dir, file := filepath.Split(src.File)
			return fmt.Sprintf("%s:%d", filepath.Join(filepath.Base(dir), file), src.Line)
		}
	}
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: opts.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				// Remove time from the output as journald will add its own timestamp and
				// we want the level first for journald marker to be effective
				return slog.Attr{}
			case slog.LevelKey:
				// Customize the name of the level key for pretty printing and the output string,
				// including custom level values
				level := a.Value.Any().(slog.Level)
				switch {
				case level < LevelInfo:
					a.Key = prefixDebugStr
					a.Value = slog.StringValue(str(LevelDebugStr, level-LevelDebug))
				case level < LevelNotice:
					a.Key = prefixInfoStr
					a.Value = slog.StringValue(str(LevelInfoStr, level-LevelInfo))
				case level < LevelWarning:
					a.Key = prefixNoticeStr
					a.Value = slog.StringValue(str(LevelNoticeStr, level-LevelNotice))
				case level < LevelError:
					a.Key = prefixWarningStr
					a.Value = slog.StringValue(str(LevelWarningStr, level-LevelWarning))
				case level < LevelCritical:
					a.Key = prefixErrorStr
					a.Value = slog.StringValue(str(LevelErrorStr, level-LevelError))
				case level < LevelAlert:
					a.Key = prefixCriticalStr
					a.Value = slog.StringValue(str(LevelCriticalStr, level-LevelCritical))
				case level < LevelEmergency:
					a.Key = prefixAlertStr
					a.Value = slog.StringValue(str(LevelAlertStr, level-LevelAlert))
				default:
					a.Key = prefixEmergencyStr
					a.Value = slog.StringValue(str(LevelEmergencyStr, level-LevelEmergency))
				}
			case slog.SourceKey:
				a.Value = slog.StringValue(opts.SourceFormat(a.Value.Any().(*slog.Source)))
			}
			// This key does not need modification, return it as is.
			return a
		},
	})
}

func str(base string, val slog.Level) string {
	if val == 0 {
		return base
	}
	return fmt.Sprintf("%s%+d", base, val)
}
