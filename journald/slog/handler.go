package sysdjournaldslog

import (
	"log/slog"
	"os"
	"strings"

	sysdjournald "github.com/iguanesolutions/go-systemd/v5/journald"
)

const (
	// https://wiki.archlinux.org/title/Systemd/Journal#Priority_level
	LevelDebug       = slog.LevelDebug
	LevelDebugStr    = "DEBUG"
	LevelInfo        = slog.LevelInfo
	LevelInfoStr     = "INFO"
	LevelNotice      = slog.Level(2)
	LevelNoticeStr   = "NOTICE"
	LevelWarning     = slog.LevelWarn
	LevelWarningStr  = "WARNING"
	LevelError       = slog.LevelError
	LevelErrorStr    = "ERROR"
	LevelCritical    = slog.Level(10)
	LevelCriticalStr = "CRITICAL"
	LevelAlert       = slog.Level(12)
	LevelAlertStr    = "ALERT"
	// LevelEmergency should not be used by applications
	LevelEmergency    = slog.Level(14)
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

var (
	levelDebugValue     = slog.StringValue(LevelDebugStr)
	levelInfoValue      = slog.StringValue(LevelInfoStr)
	levelNoticeValue    = slog.StringValue(LevelNoticeStr)
	levelWarningValue   = slog.StringValue(LevelWarningStr)
	levelErrorValue     = slog.StringValue(LevelErrorStr)
	levelCriticalValue  = slog.StringValue(LevelCriticalStr)
	levelAlertValue     = slog.StringValue(LevelAlertStr)
	levelEmergencyValue = slog.StringValue(LevelEmergencyStr)
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

// NewHandler returns a new slog handler that writes logs in a journald compatible/enhanced format.
func NewHandler(logLevel slog.Leveler) slog.Handler {
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
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
					a.Value = levelDebugValue
				case level < LevelNotice:
					a.Key = prefixInfoStr
					a.Value = levelInfoValue
				case level < LevelWarning:
					a.Key = prefixNoticeStr
					a.Value = levelNoticeValue
				case level < LevelError:
					a.Key = prefixWarningStr
					a.Value = levelWarningValue
				case level < LevelCritical:
					a.Key = prefixErrorStr
					a.Value = levelErrorValue
				case level < LevelAlert:
					a.Key = prefixCriticalStr
					a.Value = levelCriticalValue
				case level < LevelEmergency:
					a.Key = prefixAlertStr
					a.Value = levelAlertValue
				default:
					a.Key = prefixEmergencyStr
					a.Value = levelEmergencyValue
				}
			}
			// This key does not need modification, return it as is.
			return a
		},
	})
}
