package db

// Settings 配置
type Settings struct {
	Name  string // 名称
	Key   string // 键
	Value string // 值
}

// Account 账户
type Account struct {
	PlatformKey  string // 平台键值(Key为平台的唯一标识)
	PlatformName string // 平台名称(平台别名)
	Username     string // 用户名(每个平台用户名唯一)
	Password     string // 密码
	Cookies      string // cookies：保存为字符串
}

// Platform 平台
type Platform struct {
	Key    string // 唯一标识
	Name   string // 别名
	Config string // 平台配置JSON
}
