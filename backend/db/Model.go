package db

import (
	"fmt"
	"strconv"

	"github.com/go-rod/rod/lib/proto"
	"gorm.io/gorm"
)

// Setting 配置
type Setting struct {
	Layer    string // 层级
	Key      string `gorm:"primaryKey"` // 键
	Alias    string // 别名
	Category string // 分类
	Value    string // 值
}

// Account 账户
type Account struct {
	ID            uint                        `gorm:"primaryKey"`    // 主键ID
	Disabled      bool                        `gorm:"default:false"` // 是否禁用
	PlatformKey   string                      // 平台键值(Key为平台的唯一标识)
	PlatformAlias string                      // 平台名称(平台别名)
	LoginType     string                      // 登录类型(如：PWD, QRCODE)
	Username      string                      // 用户名(每个平台用户名唯一)
	Password      string                      // 密码
	Cookies       []*proto.NetworkCookieParam `gorm:"serializer:json"` // cookies：保存为字符串
}

// Platform 平台
type Platform struct {
	Key     string                 `gorm:"primaryKey"` // 唯一标识
	Alias   string                 // 别名
	Version string                 // 版本
	Config  map[string]interface{} `gorm:"serializer:json"` // 平台配置JSON
}

// BeforeCreate 创建account时，查看是同一platformKey中是否有重复的username，如果有，则自动修改username为数字累计并查看修改后的username是否重复，直到不重复为止
func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	// 创建account时，查看是同一platformKey中是否有重复的username，如果有，则自动修改username为数字累计并查看修改后的username是否重复，直到不重复为止
	var count int64
	tx.Model(Account{}).Where("platform_key = ?", account.PlatformKey).Where("username = ?", account.Username).Count(&count)
	if count > 0 {
		var newUsername string
		for i := 1; ; i++ {
			newUsername = account.Username + "_" + strconv.Itoa(i)
			tx.Model(Account{}).Where("platform_key = ?", account.PlatformKey).Where("username = ?", newUsername).Count(&count)
			if count == 0 {
				account.Username = newUsername
				break
			}
		}
	}

	return
}

// BeforeUpdate 更新account时，查看是同一platformKey中是否有重复的username，则返回错误
func (account *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	// 更新account时，查看是同一platformKey中是否有重复的username，则返回错误
	fmt.Println("account:", account)
	var count int64
	tx.Model(Account{}).Where("id <> ?", account.ID).Where("platform_key = ?", account.PlatformKey).Where("username = ?", account.Username).Count(&count)
	fmt.Println("count:", count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	return
}
