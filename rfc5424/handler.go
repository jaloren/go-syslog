package rfc5424

import (
	"fmt"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	sdId = `mdc@1806`
)

var _ slog.Handler = &Handler{}

func NewHandler(level slog.Level, appName string, w io.Writer) slog.Handler {
	var (
		err      error
		hostname string
	)
	if hostname, err = os.Hostname(); err != nil {
		hostname = "localhost"
	}
	handler := &Handler{
		level:    level,
		hostname: hostname,
		appName:  appName,
		w:        w,
	}
	return handler
}

type Handler struct {
	level    slog.Level
	hostname string
	appName  string
	w        io.Writer
	attrs    []slog.Attr
	groups   []string
	mu       sync.Mutex
}

func (s *Handler) Enabled(level slog.Level) bool {
	return s.level <= level
}

func (s *Handler) Handle(r slog.Record) error {
	if !s.Enabled(r.Level) {
		return nil
	}
	sysLogMsg := newMsg()
	file, line := r.SourceLine()
	sysLogMsg.Timestamp(r.Time.UTC()).
		AppName(s.appName).
		Hostname(s.hostname).
		Message(r.Message).
		SdParam(sdId, "log_level", r.Level.String()).
		SdParam(sdId, "source", fmt.Sprintf("%s:%d", file, line))

	// TODO: handle default case
	switch r.Level {
	case slog.ErrorLevel:
		sysLogMsg.Severity(ErrorSeverity)
	case slog.WarnLevel:
		sysLogMsg.Severity(WarningSeverity)
	case slog.DebugLevel:
		sysLogMsg.Severity(DebugSeverity)
	case slog.InfoLevel:
		sysLogMsg.Severity(InfoSeverity)
	}

	for _, attr := range s.attrs {
		s.appendAttr(attr, sysLogMsg)
	}
	r.Attrs(func(attr slog.Attr) {
		s.appendAttr(attr, sysLogMsg)
	})
	data := sysLogMsg.Build()
	data = append(data, "\n"...)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := s.w.Write(data); err != nil {
		return err
	}
	return nil
}

func (s *Handler) appendAttr(a slog.Attr, msg *syslogMsg) {
	if a.Key == "" {
		return
	}
	v := a.Value.Resolve()
	if v.Kind() == slog.GroupKind {
		for _, member := range v.Group() {
			key := s.createKey(fmt.Sprintf("%s.%s", a.Key, member.Key))
			msg.SdParam(sdId, key, member.Value.String())
		}
	} else {
		msg.SdParam(sdId, s.createKey(a.Key), a.Value.String())
	}
}

func (s *Handler) createKey(key string) string {
	if len(s.groups) == 0 {
		return key
	}
	return fmt.Sprintf("%s.%s", strings.Join(s.groups, "."), key)
}

func (s *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return s
	}
	newHandle := &Handler{
		level:    s.level,
		hostname: s.hostname,
		appName:  s.appName,
		groups:   s.groups,
		w:        s.w,
	}
	if len(s.attrs) > 0 {
		newHandle.attrs = append(newHandle.attrs, s.attrs...)
	}
	newHandle.attrs = append(newHandle.attrs, attrs...)
	return newHandle
}

func (s *Handler) WithGroup(name string) slog.Handler {
	newHandle := &Handler{
		level:    s.level,
		hostname: s.hostname,
		appName:  s.appName,
		attrs:    s.attrs,
		w:        s.w,
	}
	if len(s.groups) > 0 {
		newHandle.groups = append(newHandle.groups, s.groups...)
	}
	newHandle.groups = append(newHandle.groups, name)
	return newHandle
}
