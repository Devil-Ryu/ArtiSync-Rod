package controller

import (
	"ArtiSync-Rod/backend/db"
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBController 数据库控制器
type DBController struct {
	Ctx context.Context
	DB  *gorm.DB
}

// NewDBController 构造方法
func NewDBController() *DBController {
	return &DBController{}
}

// SetContext 设置上下文并连接数据库
func (d *DBController) SetContext(ctx context.Context) {
	d.Ctx = ctx
}

// Connect 数据库连接host
func (d *DBController) Connect(host string) (err error) {

	d.DB, err = gorm.Open(sqlite.Open(host), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("DBController Connect: %w", err)
	}
	err = d.AutoMigrate()
	if err != nil {
		return fmt.Errorf("DBController Migrate: %w", err)
	}
	return nil
}

// CheckConnect 检查数据库链接
func (d *DBController) CheckConnect() (err error) {
	if d.DB == nil {
		return errors.New("DBController not connect")
	}
	return nil
}

// AutoMigrate 迁移数据库
func (d *DBController) AutoMigrate() (err error) {
	// 检查数据库连接
	err = d.CheckConnect()
	if err != nil {
		return err
	}

	// 数据库迁移
	err = d.DB.AutoMigrate(&db.Account{}, db.Settings{}, db.Platform{})
	if err != nil {
		return err
	}
	return nil
}
