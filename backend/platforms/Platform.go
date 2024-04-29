package platforms

import (
	"ArtiSync-Rod/backend/controller"
	"ArtiSync-Rod/backend/db"
	"ArtiSync-Rod/backend/utils"
	"context"
	"fmt"
	"reflect"

	"github.com/go-rod/rod/lib/proto"
	"github.com/mitchellh/mapstructure"
)

// Model 平台模型
type Model struct {
	Key             string                    // 平台唯一标识
	Alias           string                    // 平台别名
	Account         *db.Account               // 平台账号
	AccountProgress *utils.AccountInfo        // 平台账号进度
	Article         *utils.Article            // 待上传文章
	Ctx             context.Context           // 上下文
	RODController   *controller.RODController // 机器人控制器
	DBController    *controller.DBController  // 数据库控制器
}

// SetArticle 设置待上传文章
func (m *Model) SetArticle(article *utils.Article) {
	m.Article = article
}

// SetAccount 设置平台账号
func (m *Model) SetAccount(account *db.Account, accountProgress *utils.AccountInfo) {
	m.Account = account
	m.AccountProgress = accountProgress
}

// SetController 设置控制器(必要-设置平台控制器)
func (m *Model) SetController(dbController *controller.DBController, rodController *controller.RODController) {
	m.DBController = dbController
	m.RODController = rodController
}

// HasDBController 是否有数据库控制器
func (m *Model) HasDBController() bool {
	return m.DBController != nil
}

// HasRODController 是否有机器人控制器
func (m *Model) HasRODController() bool {
	return m.RODController != nil
}

// HasAccount 是否有平台账号
func (m *Model) HasAccount() bool {
	return m.Account != nil
}

// Start 启动平台(必要-需要重写)
func (m *Model) Start(ctx context.Context, dbc *controller.DBController, rdc *controller.RODController, config interface{}, account *db.Account, article *utils.Article, accountProgress *utils.AccountInfo, publishFunc func() error) (err error) {
	err = m.InitRod(dbc, rdc, &config) // 初始化机器人
	if err != nil {
		return err
	}
	m.Ctx = ctx                            // 设置上下文
	m.SetAccount(account, accountProgress) // 设置账号信息
	m.SetArticle(article)                  // 设置文章信息
	err = publishFunc()                    // 发布文章
	if err != nil {
		return err
	}

	return nil
}

// InitRod 初始化Rod
func (m *Model) InitRod(dbc *controller.DBController, rdc *controller.RODController, config interface{}) (err error) {
	// 设置controller
	m.SetController(dbc, rdc)

	// 加载配置
	configMap, err := m.LoadConfig()
	if err != nil {
		return err
	}

	// 将配置读取到结构体中
	err = mapstructure.Decode(configMap, config)
	if err != nil {
		return fmt.Errorf("配置解码失败: %w", err)
	}

	return err
}

// LoadConfig 从数据库中加载配置(必要-加载平台配置)
func (m *Model) LoadConfig() (config map[string]interface{}, err error) {
	if m.HasDBController() == false { // 如果没有数据库控制器则返回错误
		return config, fmt.Errorf("LoadConfig: 没有数据库控制器")
	}
	platformInfo, err := m.DBController.GetPlatform(db.Platform{Key: m.Key}) // 获取平台信息
	if err != nil {                                                          // 如果获取失败则返回错误
		return config, err
	}

	config = platformInfo.Config // 读取配置
	return config, err
}

// OpenPage 打开页面(必要-统一方法)
func (m *Model) OpenPage(pageURL string, accountID uint) (err error) {
	// 检查配置
	if m.HasDBController() == false { // 如果没有数据库控制器则返回错误
		return fmt.Errorf("OpenPage: 没有数据库控制器")
	}

	account, err := m.DBController.GetAccount(db.Account{ID: accountID}) // 获取账号信息
	if err != nil {
		return fmt.Errorf("OpenPage: %w", err)
	}

	m.SetAccount(&account, nil)

	// 获取平台账号的cookies
	if m.HasAccount() == false { // 如果没有平台账号则返回错误
		return fmt.Errorf("OpenPage: 没有平台账号")
	}

	/*设置浏览器*/
	rdc := controller.NewRODController()

	rdc.StartBrowser(false)  // 显示浏览器
	defer rdc.CloseBrowser() // 关闭浏览器

	/*设置浏览器cookies*/
	err = rdc.Browser.SetCookies(m.Account.Cookies)
	if err != nil {
		return fmt.Errorf("SetCookies: %w", err)
	}

	page := rdc.Browser.MustPage(pageURL)

	// 等待窗口关闭
	wait := page.EachEvent(func(e *proto.PageLifecycleEvent) (stop bool) {
		return true
	})
	wait()

	return nil
}

// CheckConfig 检查配置是否正确(TODOisNil不生效)
func (m *Model) CheckConfig(config interface{}) (err error) {
	if isNil(config) {
		return fmt.Errorf(m.Key + "配置未设置")
	}

	// 判断DBController是否设置
	if m.HasDBController() == false {
		return fmt.Errorf(m.Key + " DBController未设置")
	}

	// 判断RODController是否设置
	if m.HasRODController() == false {
		return fmt.Errorf(m.Key + " RODController未设置")
	}

	return nil
}

// CheckAuthentication 检查授权(必要-需要重写)
func (m *Model) CheckAuthentication() (authInfo map[string]string, err error) { return authInfo, err }

// Login 登录(必要-需要重写)
func (m *Model) Login() (accounts []db.Account, err error) { return accounts, err }

// Publish 平台运行方法(必要-需要重写)
func (m *Model) Publish() (err error) { return err }

// UpdatePlatformInfo 更新平台信息(必要-需要重写)
func (m *Model) UpdatePlatformInfo() {}

// 判断是否为空
func isNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
