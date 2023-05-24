package repository

import (
	"github.com/WebEngrChild/go-graphql-server/pkg/domain/model"
)

type Message interface {
	GetMessages() ([]*model.Message, error)
	CreateMessage(todos *model.Message) (*model.Message, error)
}
