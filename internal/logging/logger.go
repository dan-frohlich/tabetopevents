package logging

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/dan-frohlich/tabetopevents/internal/tui"
)

type Logger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
	Fatal(message string, args ...any)
	WithError(error) Logger
	WithField(key string, value any) Logger
	WithFields(fields ...any) Logger
}

var (
	_ Logger = Log{}
	// _ Logger = LogEvent{}
)

type Log struct {
	Level LogLevel
}

type LogStyle struct {
	Icon  rune
	Style lipgloss.Style
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

func (ll LogLevel) LogStyle() LogStyle {
	switch ll {
	case LogLevelDebug:
		return DebugStyle
	case LogLevelInfo:
		return InfoStyle
	case LogLevelWarn:
		return WarnStyle
	case LogLevelError:
		return ErrorStyle
	case LogLevelFatal:
		return FatalStyle
	default:
		return LogStyle{Icon: '?', Style: DebugStyle.Style}
	}
}

var (
	DebugStyle = LogStyle{
		Icon:  'd',
		Style: lipgloss.NewStyle().Foreground(tui.ColorGray).Italic(true),
	}
	InfoStyle = LogStyle{
		Icon:  'i',
		Style: lipgloss.NewStyle().Foreground(tui.ColorBlue).Italic(true),
	}
	WarnStyle = LogStyle{
		Icon:  'w',
		Style: lipgloss.NewStyle().Foreground(tui.ColorYellow).Italic(true),
	}
	ErrorStyle = LogStyle{
		Icon:  'e',
		Style: lipgloss.NewStyle().Foreground(tui.ColorRed).Italic(true),
	}
	FatalStyle = LogStyle{
		Icon:  'f',
		Style: lipgloss.NewStyle().Foreground(tui.ColorRed).Bold(true).Italic(true),
	}
)

func formatLogMessage(ll LogLevel, message string, args ...any) string {
	s := fmt.Sprintf("[%c] %s", ll.LogStyle().Icon, message)
	if len(args) == 0 {
		return ll.LogStyle().Style.Render(s)
	}
	if len(args)%2 != 0 {
		var az []any
		az = append(az, ll.LogStyle().Style.Render(s))
		az = append(az, "?odd-arg-size?")
		az = append(az, args...)
		return strings.Join(argsToStrings(az), " ")
	}
	var sz []string
	sz = append(sz, s)
	for i := range args {
		if i%2 == 0 {
			continue
		}
		k := argToString(args[i-1])
		v := argToString(args[i])
		sz = append(sz, fmt.Sprintf("%s=%s", k, v))
	}
	return ll.LogStyle().Style.Render(strings.Join(sz, " "))
}

// Debug implements Logger.
func (l Log) Debug(message string, args ...any) {
	var ll LogLevel = LogLevelDebug
	if l.Level > ll {
		return
	}
	fmt.Println(formatLogMessage(ll, message, args...))
}

func argsToStrings(args []any) (sz []string) {
	sz = make([]string, 0, len(args))
	for _, arg := range args {
		sz = append(sz, argToString(arg))
	}
	return sz
}

func argToString(arg any) string {
	switch v := arg.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	case int, int16, int32, int64, int8, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Error implements Logger.
func (l Log) Error(message string, args ...any) {
	var ll LogLevel = LogLevelError
	if l.Level > ll {
		return
	}
	fmt.Println(formatLogMessage(ll, message, args...))
}

// Fatal implements Logger.
func (l Log) Fatal(message string, args ...any) {
	var ll LogLevel = LogLevelFatal
	if l.Level > ll {
		return
	}
	fmt.Println(formatLogMessage(ll, message, args...))
}

// Info implements Logger.
func (l Log) Info(message string, args ...any) {
	var ll LogLevel = LogLevelInfo
	if l.Level > ll {
		return
	}
	fmt.Println(formatLogMessage(ll, message, args...))
}

// Warn implements Logger.
func (l Log) Warn(message string, args ...any) {
	var ll LogLevel = LogLevelWarn
	if l.Level > ll {
		return
	}
	fmt.Println(formatLogMessage(ll, message, args...))
}

// WithError implements Logger.
func (l Log) WithError(error) Logger {
	panic("unimplemented")
}

// WithField implements Logger.
func (l Log) WithField(key string, value any) Logger {
	panic("unimplemented")
}

// WithFields implements Logger.
func (l Log) WithFields(fields ...any) Logger {
	panic("unimplemented")
}

// type LogEvent struct {
// 	fields map[string]any
// 	err    error
// 	Level  LogLevel
// }

// // Debug implements Logger.
// func (l LogEvent) Debug(message string, args ...any) {
// 	var ll LogLevel = LogLevelDebug
// 	if l.Level > ll {
// 		return
// 	}
// 	fmt.Println(formatLogMessage(ll, message, args...))
// }

// // Error implements Logger.
// func (l LogEvent) Error(message string, args ...any) {
// 	var ll LogLevel = LogLevelError
// 	if l.Level > ll {
// 		return
// 	}
// 	fmt.Println(formatLogMessage(ll, message, args...))
// }

// // Fatal implements Logger.
// func (l LogEvent) Fatal(message string, args ...any) {
// 	var ll LogLevel = LogLevelFatal
// 	if l.Level > ll {
// 		return
// 	}
// 	fmt.Println(formatLogMessage(ll, message, args...))
// }

// // Info implements Logger.
// func (l LogEvent) Info(message string, args ...any) {
// 	var ll LogLevel = LogLevelInfo
// 	if l.Level > ll {
// 		return
// 	}
// 	fmt.Println(formatLogMessage(ll, message, args...))
// }

// // Warn implements Logger.
// func (l LogEvent) Warn(message string, args ...any) {
// 	var ll LogLevel = LogLevelWarn
// 	if l.Level > ll {
// 		return
// 	}
// 	fmt.Println(formatLogMessage(ll, message, args...))
// }

// // WithError implements Logger.
// func (l LogEvent) WithError(error) Logger {
// 	panic("unimplemented")
// }

// // WithField implements Logger.
// func (l LogEvent) WithField(key string, value any) Logger {
// 	panic("unimplemented")
// }

// // WithFields implements Logger.
// func (l LogEvent) WithFields(fields ...any) Logger {
// 	panic("unimplemented")
// }
