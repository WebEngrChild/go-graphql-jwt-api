package repository

import (
	"context"

	"github.com/WebEngrChild/go-graphql-server/pkg/domain/model/graph"
)

type User interface {
	GetMapInIDs(ctx context.Context, ids []string) (map[string]*graph.User, error)
	GetUserById(id string) (*graph.User, error)
	GetUserByEmail(email string) (*graph.User, error)
	Encrypt(plain string) (string, error)
	Decrypt(encrypted string) (string, error)
}
