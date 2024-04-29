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
	err = d.DB.AutoMigrate(&db.Account{}, db.Setting{}, db.Platform{})
	if err != nil {
		return err
	}
	return nil
}

// CreateSettings 创建配置
func (d *DBController) CreateSettings(settings []db.Setting) ([]db.Setting, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return settings, err
	}

	result := d.DB.Create(&settings)
	if result.Error != nil {
		return settings, result.Error
	}
	return settings, nil
}

// CreateAccounts 创建账户
func (d *DBController) CreateAccounts(accounts []db.Account) ([]db.Account, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return accounts, err
	}

	result := d.DB.Create(&accounts)
	if result.Error != nil {
		return accounts, result.Error
	}
	return accounts, nil
}

// CreatePlatforms 创建平台
func (d *DBController) CreatePlatforms(platform []db.Platform) ([]db.Platform, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return platform, err
	}

	result := d.DB.Create(&platform)
	if result.Error != nil {
		return platform, result.Error
	}
	return platform, nil
}

// UpdateSetting 更新配置
func (d *DBController) UpdateSetting(setting db.Setting) error {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return err
	}

	result := d.DB.Save(&setting)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateAccount 更新账户
func (d *DBController) UpdateAccount(account db.Account) (db.Account, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return account, err
	}

	result := d.DB.Save(&account)
	if result.Error != nil {
		return account, result.Error
	}
	return account, nil
}

// UpdatePlatform 更新平台
func (d *DBController) UpdatePlatform(platform db.Platform) error {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return err
	}

	result := d.DB.Save(&platform)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// QuerySettings 通过查询条件获取配置
func (d *DBController) QuerySettings(query map[string]interface{}) ([]db.Setting, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return nil, err
	}

	var settings []db.Setting
	result := d.DB.Where(query).Find(&settings)
	if result.Error != nil {
		return nil, result.Error
	}

	return settings, nil
}

// QueryAccounts 通过查询条件获取账户
func (d *DBController) QueryAccounts(query map[string]interface{}) ([]db.Account, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return nil, err
	}

	var accounts []db.Account
	result := d.DB.Where(query).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

// QueryPlatforms 通过查询条件获取平台
func (d *DBController) QueryPlatforms(query map[string]interface{}) ([]db.Platform, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return nil, err
	}

	var platforms []db.Platform
	result := d.DB.Where(query).Find(&platforms)
	if result.Error != nil {
		return nil, result.Error
	}
	return platforms, nil
}

// GetSetting 获取单个配置
func (d *DBController) GetSetting(setting db.Setting) (db.Setting, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return setting, err
	}

	result := d.DB.Model(&db.Setting{}).First(&setting)

	if result.Error != nil {
		return setting, result.Error
	}
	return setting, nil
}

// GetAccount 获取单个账户
func (d *DBController) GetAccount(account db.Account) (db.Account, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return account, err
	}

	result := d.DB.First(&account, account.ID)
	if result.Error != nil {
		return account, result.Error
	}

	return account, nil
}

// GetPlatform 获取单个平台
func (d *DBController) GetPlatform(platform db.Platform) (db.Platform, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return platform, err
	}

	result := d.DB.Model(&db.Platform{}).First(&platform)
	if result.Error != nil {
		return platform, result.Error
	}
	return platform, nil
}

// DeleteSetting 删除配置
func (d *DBController) DeleteSetting(setting db.Setting) error {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return err
	}

	result := d.DB.Delete(&setting)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteAccount 删除账户
func (d *DBController) DeleteAccount(account db.Account) error {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return err
	}

	result := d.DB.Delete(&account)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeletePlatform 删除平台
func (d *DBController) DeletePlatform(platform db.Platform) error {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return err
	}

	result := d.DB.Delete(&platform)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
