package log

import (
	"github.com/jaloren/go-syslog/rfc5424"
	"golang.org/x/exp/slog"
	"io"
	"testing"
	"time"
)

type Name struct {
	First, Last string
}

func (n Name) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("first", n.First),
		slog.String("last", n.Last))
}

func BenchmarkFluentLogger(b *testing.B) {
	b.ReportAllocs()
	logger := New(slog.InfoLevel, "demo", io.Discard)
	for n := 0; n < b.N; n++ {
		logger.Log("this is a message").
			String("account", "bank").
			Time("created_at", time.Now().UTC()).
			Duration("ttl", time.Second*4).
			Int("count", 1).
			Bool("access", true).
			Any("person", Name{First: "peter", Last: "parker"}).
			Warn()
	}
}

func BenchmarkStdLogger(b *testing.B) {
	b.ReportAllocs()
	logger := slog.New(rfc5424.NewHandler(slog.WarnLevel, "demo", io.Discard))
	for n := 0; n < b.N; n++ {
		attrs := []slog.Attr{
			slog.String("account", "bank"),
			slog.Time("created_at", time.Now().UTC()),
			slog.Duration("ttl", time.Second*4),
			slog.Int("count", 1),
			slog.Bool("access", true),
			slog.Any("person", Name{First: "peter", Last: "parker"}),
		}
		logger.LogAttrs(slog.WarnLevel, "this is a message", attrs...)
	}
}
