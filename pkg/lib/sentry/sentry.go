package sentry

import (
	"strings"

	"github.com/WebEngrChild/go-graphql-server/pkg/lib/config"
	"github.com/getsentry/sentry-go"
	"golang.org/x/xerrors"
)

func SetUp(c *config.Config) error {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              c.SentryDsn,
		Environment:      c.Env,
		Debug:            true,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
		// BeforeSend のフックで Event を書き換え
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			for i := range event.Exception {
				exception := &event.Exception[i]
				// fmt.wrapError, xerrors.wrapError 以外は何もしない
				if !strings.Contains(exception.Type, "wrapError") {
					continue
				}
				// 最初の : で分割(正しく Wrapされていないものは無視)
				sp := strings.SplitN(exception.Value, ":", 2)
				if len(sp) != 2 {
					continue
				}
				// : の前を Typeに、 : より後ろを Value に
				exception.Type, exception.Value = sp[0], sp[1]
			}
			return event
		},
	}); err != nil {
		return xerrors.Errorf("fail to init sentry: %w", err)
	}

	return nil
}
