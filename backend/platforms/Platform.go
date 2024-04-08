package platforms

import (
	"ArtiSync-Rod/backend/controller"
	"ArtiSync-Rod/backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/go-rod/rod/lib/proto"
)

// Model 平台模型
type Model struct {
	Article       *utils.Article              // 待上传文章
	Cookies       []*proto.NetworkCookieParam // cookies
	Ctx           context.Context             // 上下文
	Name          string                      // 名称
	PlatformIndex int                         // 平台序号
	RODController *controller.RODController   // 机器人控制器
}

// Init 初始化平台(统一方法)
func (m *Model) Init(ctx context.Context, article *utils.Article, platformIndex int) {
	m.Article = article
	m.PlatformIndex = platformIndex
	m.Ctx = ctx
}

// SetCookies 设置cookies(统一方法)
func (m *Model) SetCookies(Cookies []*proto.NetworkCookieParam) {
	m.Cookies = Cookies
}

// LoadCookies 读取cookies(统一方法)
func (m *Model) LoadCookies() (cookies []*proto.NetworkCookieParam, err error) {
	cookiePath, err := m.RODController.GetCookiePath(m.Name)
	if err != nil {
		return cookies, err
	}
	dataBytes, err := os.ReadFile(cookiePath)
	if err != nil {
		return cookies, fmt.Errorf("LoadCookies: %w", err)
	}

	// 读取文件内容
	err = json.Unmarshal(dataBytes, &cookies)
	if err != nil {
		return cookies, fmt.Errorf("LoadCookies: %w", err)
	}
	m.Cookies = cookies
	return cookies, nil
}

// LoadConfig 加载配置(必要-加载平台配置)
func (m *Model) LoadConfig(defaultConfig map[string]interface{}, forceDefault bool) (config map[string]interface{}, err error) {
	cutl := utils.NewCommonUtils()
	configPath, err := m.RODController.GetPlatformConfigPath(m.Name)
	if err != nil {
		return config, err
	}

	// 尝试加载配置文件
	config, err = cutl.LoadJSONFile(configPath)
	if err != nil || forceDefault { // 如果不存在配置文件则写入默认配置
		log.Println("加载默认配置")
		err = cutl.SaveJSONFile(configPath, defaultConfig)
		if err != nil {
			return config, err
		}
	}

	// 读取配置
	config, err = cutl.LoadJSONFile(configPath)

	if err != nil {
		return config, err
	}

	return config, nil
}

// CheckConfig 检查config(统一方法)
func (m *Model) CheckConfig(config interface{}) (err error) {
	// 进行配置校验
	validate := validator.New()
	err = validate.Struct(config) // 校验config必要字段是否已有
	if err != nil {
		log.Println("配置校验失败: ", err)
		return err
	}

	log.Println("配置校验通过")

	return nil
}

// OpenPage 打开页面(统一方法)
func (m *Model) OpenPage(pageURL string) (err error) {
	// 首先获取cookies
	_, err = m.LoadCookies()
	if err != nil {
		return err
	}

	/*设置浏览器*/
	if m.RODController.Browser == nil {
		m.RODController.StartBrowser(false) // 启动浏览器
	}

	/*设置浏览器cookies*/
	m.RODController.Browser.SetCookies(m.Cookies)

	page := m.RODController.Browser.MustPage(pageURL)

	// 等待窗口关闭
	wait := page.EachEvent(func(e *proto.PageLifecycleEvent) (stop bool) {
		return true
	})
	wait()

	return nil
}

// SetConfig 加载配置(必要-需要重写)
func (m *Model) SetConfig(foreDefault bool) (err error) { return err }

// CheckAuthentication 检查授权(必要-需要重写)
func (m *Model) CheckAuthentication() (authInfo map[string]string, err error) { return authInfo, err }

// Login 登录(必要-需要重写)
func (m *Model) Login() (err error) { return err }

// Run 平台运行方法(必要-需要重写)
func (m *Model) Run() (err error) { return err }
