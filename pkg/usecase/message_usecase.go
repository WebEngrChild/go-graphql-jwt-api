package usecase

import (
	"fmt"
	"time"

	"github.com/WebEngrChild/go-graphql-server/pkg/domain/model"
	"github.com/WebEngrChild/go-graphql-server/pkg/domain/repository"
)

type Message interface {
	CreateMessage(text *string, userId *string) (*model.Message, error)
	GetMessages() ([]*model.Message, error)
}

type MsgUseCase struct {
	msgRepo repository.Message
}

func NewMsgUseCase(msgRepo repository.Message) Message {
	MessageUseCase := MsgUseCase{msgRepo: msgRepo}
	return &MessageUseCase
}

func (m *MsgUseCase) CreateMessage(text *string, userId *string) (*model.Message, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	todo := model.Message{
		Text:      *text,
		UserID:    *userId,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	created, err := m.msgRepo.CreateMessage(&todo)
	if err != nil {
		return nil, fmt.Errorf("useCase CreateMessage() err: %w", err)
	}

	return created, nil
}

func (m *MsgUseCase) GetMessages() ([]*model.Message, error) {
	todos, err := m.msgRepo.GetMessages()
	if err != nil {
		return nil, fmt.Errorf("resolver Todos() err %w", err)
	}

	return todos, nil
}
