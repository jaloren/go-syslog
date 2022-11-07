package log

import (
	"github.com/jaloren/go-syslog/rfc5424"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"time"
)

func Default(appName string) Logger {
	return New(slog.InfoLevel, appName, os.Stdout)
}

func New(level slog.Level, appName string, w io.Writer) Logger {
	handler := rfc5424.NewHandler(level, appName, w)
	return Logger{
		std: slog.New(handler),
	}
}

type Logger struct {
	std slog.Logger
}

func (l Logger) Log(msg string) *Message {
	return &Message{msg: msg, logger: l}
}

type Message struct {
	logger Logger
	msg    string
	attrs  []slog.Attr
}

func (m *Message) String(key, val string) *Message {
	m.attrs = append(m.attrs, slog.String(key, val))
	return m
}

func (m *Message) Int(key string, val int) *Message {
	m.attrs = append(m.attrs, slog.Int(key, val))
	return m
}

func (m *Message) Int64(key string, val int64) *Message {
	m.attrs = append(m.attrs, slog.Int64(key, val))
	return m
}

func (m *Message) Time(key string, val time.Time) *Message {
	m.attrs = append(m.attrs, slog.Time(key, val))
	return m
}

func (m *Message) Duration(key string, val time.Duration) *Message {
	m.attrs = append(m.attrs, slog.Duration(key, val))
	return m
}

func (m *Message) Float64(key string, val float64) *Message {
	m.attrs = append(m.attrs, slog.Float64(key, val))
	return m
}

func (m *Message) Bool(key string, val bool) *Message {
	m.attrs = append(m.attrs, slog.Bool(key, val))
	return m
}

func (m *Message) Group(key string, as ...slog.Attr) *Message {
	m.attrs = append(m.attrs, slog.Group(key, as...))
	return m
}

func (m *Message) Any(key string, val any) *Message {
	m.attrs = append(m.attrs, slog.Any(key, val))
	return m
}

func (m *Message) Info() {
	m.logger.std.LogAttrs(slog.InfoLevel, m.msg, m.attrs...)
}

func (m *Message) Error(err error) {
	if err != nil {
		m.Any("err", err)
	}
	m.logger.std.LogAttrs(slog.ErrorLevel, m.msg, m.attrs...)
}

func (m *Message) Warn() {
	m.logger.std.LogAttrs(slog.WarnLevel, m.msg, m.attrs...)
}

func (m *Message) Debug() {
	m.logger.std.LogAttrs(slog.DebugLevel, m.msg, m.attrs...)
}
