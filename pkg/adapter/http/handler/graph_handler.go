package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/WebEngrChild/go-graphql-server/pkg/adapter/http/resolver"
	"github.com/WebEngrChild/go-graphql-server/pkg/lib/graph/generated"
	"github.com/WebEngrChild/go-graphql-server/pkg/lib/graph/loader"
	"github.com/WebEngrChild/go-graphql-server/pkg/usecase"
	"github.com/graph-gophers/dataloader"
	"github.com/labstack/echo/v4"
)

type Graph interface {
	QueryHandler() echo.HandlerFunc
}

type GraphHandler struct {
	MsgUseCase  usecase.Message
	UserUseCase usecase.User
}

func NewGraphHandler(mu usecase.Message, uc usecase.User) Graph {
	GraphHandler := GraphHandler{
		MsgUseCase:  mu,
		UserUseCase: uc,
	}
	return &GraphHandler
}

func (g *GraphHandler) QueryHandler() echo.HandlerFunc {
	ldr := &loader.Loaders{
		UserLoader: dataloader.NewBatchedLoader(
			g.UserUseCase.BatchGetUsers,
			dataloader.WithCache(&dataloader.NoCache{}),
		),
	}

	rslvr := resolver.Resolver{
		MsgUseCase:  g.MsgUseCase,
		UserUseCase: g.UserUseCase,
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &rslvr}),
	)

	return func(c echo.Context) error {
		loader.Middleware(ldr, srv).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
