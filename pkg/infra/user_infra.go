package infra

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/WebEngrChild/go-graphql-server/pkg/domain/model/graph"
	"github.com/WebEngrChild/go-graphql-server/pkg/domain/repository"
	"github.com/WebEngrChild/go-graphql-server/pkg/lib/config"
	"gorm.io/gorm"
)

type userRepository struct {
	db     *gorm.DB
	config *config.Config
}

func NewUserRepository(db *gorm.DB, config *config.Config) repository.User {
	return &userRepository{db, config}
}

func (r *userRepository) GetMapInIDs(ctx context.Context, ids []string) (map[string]*graph.User, error) {
	var users []*graph.User
	if err := r.db.Find(&users, ids).Error; err != nil {
		return nil, xerrors.Errorf("repository GetMapInIDs err %w", err)
	}
	results := make(map[string]*graph.User)
	for _, user := range users {
		results[user.ID] = user
	}

	return results, nil
}

func (r *userRepository) GetUserByEmail(email string) (*graph.User, error) {
	var user *graph.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, xerrors.Errorf("get users by email failed , %w", err)
	}

	return user, nil
}

func (r *userRepository) GetUserById(id string) (*graph.User, error) {
	var user *graph.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, xerrors.Errorf("get users by id failed , %w", err)
	}

	return user, nil
}

func (r *userRepository) Encrypt(plain string) (string, error) {
	var enc string
	err := r.db.Raw("SELECT HEX(AES_ENCRYPT(?, ?))", plain, r.config.EncryptKey).Scan(&enc).Error
	if err != nil {
		return enc, xerrors.Errorf("encrypt email failed , %w", err)
	}

	return enc, nil
}

func (r *userRepository) Decrypt(encrypted string) (string, error) {
	var dec string
	err := r.db.Raw("SELECT CONVERT(AES_DECRYPT(UNHEX(?), ?) USING utf8)", encrypted, r.config.EncryptKey).Scan(&dec).Error
	if err != nil {
		return dec, xerrors.Errorf("decrypt email failed, %w", err)
	}

	return dec, nil
}
