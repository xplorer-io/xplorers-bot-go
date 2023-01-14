package sentrylib

import (
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/xplorer-io/xplorers-bot-go/internal/ssm"
)

func InitializeSentry() error {
	sentryDsn, err := ssm.GetSsmParameter(os.Getenv("SENTRY_DSN_SSM_PATH"))
	if err != nil {
		return err
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDsn,
		TracesSampleRate: 1.0,
		Environment:      os.Getenv("ENVIRONMENT"),
		Debug:            true,
	}); err != nil {
		return err
	}

	return nil
}
