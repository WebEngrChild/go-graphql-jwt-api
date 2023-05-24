package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/getsentry/sentry-go"

	"golang.org/x/xerrors"

	"github.com/WebEngrChild/go-graphql-server/pkg/domain/model/graph"
	"github.com/WebEngrChild/go-graphql-server/pkg/domain/repository"
	"github.com/WebEngrChild/go-graphql-server/pkg/lib/graph/loader"
	"github.com/graph-gophers/dataloader"
)

type User interface {
	BatchGetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result
	LoadUser(ctx context.Context, userID string) (*graph.User, error)
}

type UserUseCase struct {
	userRepo repository.User
}

func NewUserUseCase(userRepo repository.User) User {
	UserUseCase := UserUseCase{userRepo: userRepo}
	return &UserUseCase
}

// BatchGetUsers dataloader.BatchFuncを実装したメソッド
// ユーザーの一覧を取得する
func (u *UserUseCase) BatchGetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// dataloader.Keysの型を[]stringに変換する
	userIDs := make([]string, len(keys))
	for ix, key := range keys {
		userIDs[ix] = key.String()
	}
	// 実処理
	log.Printf("BatchGetUsers(id = %s)\n", strings.Join(userIDs, ","))
	userByID, err := u.userRepo.GetMapInIDs(ctx, userIDs)
	if err != nil {
		err := fmt.Errorf("BatchGetUsers failed err: %w", err)
		sentry.CaptureException(err)
		log.Printf("%v\n", err)
		return nil
	}
	// []*model.User[]*dataloader.Resultに変換する
	output := make([]*dataloader.Result, len(keys))
	for index, userKey := range keys {
		user, ok := userByID[userKey.String()]
		if ok {
			output[index] = &dataloader.Result{Data: user, Error: nil}
		} else {
			err := fmt.Errorf("BatchGetUsers failed user not found %s", userKey.String())
			output[index] = &dataloader.Result{Data: nil, Error: err}
		}
	}
	return output
}

// LoadUser dataloader.Loadをwrapして型づけした実装
func (u *UserUseCase) LoadUser(ctx context.Context, userID string) (*graph.User, error) {
	log.Printf("LoadUser(id = %s)\n", userID)
	loaders := loader.GetLoaders(ctx)
	thunk := loaders.UserLoader.Load(ctx, dataloader.StringKey(userID))
	result, err := thunk()
	if err != nil {
		return nil, xerrors.Errorf("LoadUser err %w", err)
	}
	user := result.(*graph.User)
	log.Printf("return LoadUser(id = %s, name = %s)\n", user.ID, user.Name)
	return user, nil
}
