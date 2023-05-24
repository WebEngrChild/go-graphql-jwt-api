package infra

import (
	"golang.org/x/xerrors"

	"github.com/WebEngrChild/go-graphql-server/pkg/domain/model"
	"github.com/WebEngrChild/go-graphql-server/pkg/domain/repository"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repository.Message {
	return &messageRepository{db}
}

func (r *messageRepository) GetMessages() ([]*model.Message, error) {
	var records []model.Message
	if result := r.db.Find(&records); result.Error != nil {
		return nil, xerrors.Errorf("repository  GetTodos() err %w", result.Error)
	}

	var res []*model.Message
	for _, record := range records {
		record := record
		res = append(res, &record)
	}

	return res, nil
}

func (r *messageRepository) CreateMessage(todo *model.Message) (*model.Message, error) {
	if result := r.db.Create(todo); result.Error != nil {
		return nil, xerrors.Errorf("repository CreateTodo() err %w", result.Error)
	}
	return todo, nil
}
