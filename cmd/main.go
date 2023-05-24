package main

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/WebEngrChild/go-graphql-server/pkg/adapter/http/handler"
	authMiddleware "github.com/WebEngrChild/go-graphql-server/pkg/adapter/http/middleware"
	"github.com/WebEngrChild/go-graphql-server/pkg/adapter/http/route"
	"github.com/WebEngrChild/go-graphql-server/pkg/infra"
	"github.com/WebEngrChild/go-graphql-server/pkg/lib/config"
	initSentry "github.com/WebEngrChild/go-graphql-server/pkg/lib/sentry"
	"github.com/WebEngrChild/go-graphql-server/pkg/usecase"
)

func main() {
	// Config
	c, cErr := config.New()

	// Sentry
	if err := initSentry.SetUp(c); err != nil {
		sentry.CaptureException(fmt.Errorf("initSentry err: %w", err))
	}

	// DB
	db, err := infra.NewDBConnector(c)
	if err != nil {
		sentry.CaptureException(fmt.Errorf("initDb err: %w", err))
	}

	// DI
	mr := infra.NewMessageRepository(db)
	ur := infra.NewUserRepository(db, c)
	au := usecase.NewAuthUseCase(ur, c)
	mu := usecase.NewMsgUseCase(mr)
	uu := usecase.NewUserUseCase(ur)
	ch := handler.NewCsrfHandler()
	lh := handler.NewLoginHandler(au)
	gh := handler.NewGraphHandler(mu, uu)
	ph := playground.Handler("GraphQL", "/query")
	am := authMiddleware.NewAuthMiddleware(au)

	// Rooting
	r := route.NewInitRoute(ch, lh, gh, ph, am)
	_, err = r.InitRouting(c)
	if err != nil {
		sentry.CaptureException(fmt.Errorf("InitRouting at NewInitRoute err: %w", err))
	}

	defer func() {
		// .envが存在しない場合
		if cErr != nil {
			sentry.CaptureException(fmt.Errorf("config err: %w", cErr))
		}
		// panic の場合も Sentry に通知する場合は Recover() を呼ぶ
		sentry.Recover()
		// サーバーへは非同期でバッファしつつ送信するため、未送信のものを忘れずに送る(引数はタイムアウト時間)
		sentry.Flush(2 * time.Second)
	}()
}
