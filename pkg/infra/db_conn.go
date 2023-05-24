package infra

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/WebEngrChild/go-graphql-server/pkg/lib/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBConnector(cfg *config.Config) (*gorm.DB, error) {
	// TODO::main.go配下ではブレイクポイントが止まらない理由を確認
	// TODO::172.30.0.3がmysqlじゃない理由を確認
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, xerrors.Errorf("db connection failed：%w", err)
	}

	return db, nil
}
