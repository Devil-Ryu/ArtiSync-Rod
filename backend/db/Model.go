package db

import "github.com/go-rod/rod/lib/proto"

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
	Disabled      bool                        `gorm:"default:false"` // 是否禁用
	PlatformKey   string                      `gorm:"primaryKey"`    // 平台键值(Key为平台的唯一标识)
	PlatformAlias string                      // 平台名称(平台别名)
	Username      string                      `gorm:"primaryKey"` // 用户名(每个平台用户名唯一)
	LoginType     string                      // 登录类型(如：PWD, QRCODE)
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
