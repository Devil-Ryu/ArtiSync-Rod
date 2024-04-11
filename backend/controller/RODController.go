package controller

import (
	"log"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// RODController 机器人控制器
type RODController struct {
	Browser   *rod.Browser       // 浏览器
	Launcher  *launcher.Launcher // 启动器
	CheckTime int                // 检查时间
}

// NewRODController 创建一个机器人控制器
func NewRODController() *RODController {
	return &RODController{
		CheckTime: 1,
	}
}

// StartBrowser 开启一个浏览器
func (rodc *RODController) StartBrowser(headless bool) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// 设置浏览器启动器
		log.Println("正在启动浏览器")
		browserPath, ok := launcher.LookPath() // 查找浏览器路径
		if ok {                                // 如果找到本地浏览器，则使用
			rodc.Launcher = launcher.New().Bin(browserPath) // 设置浏览器
			log.Println("使用本地浏览器")
		} else { // 如果未找到则不使用本地浏览器
			rodc.Launcher = launcher.New() // 设置浏览器
			log.Println("未找到本地浏览器，使用默认浏览器")
		}
		rodc.Launcher.Headless(headless) // 设置浏览器参数
		defer rodc.Launcher.Cleanup()    // cleanup

		// 启动浏览器
		rodc.Browser = rod.New().NoDefaultDevice().ControlURL(rodc.Launcher.MustLaunch()).MustConnect() // 打开浏览器
		// 等待浏览器启动完成后打印日志显示浏览器参数
		log.Println("浏览器启动: headless->", headless)
		defer wg.Done() // 判断是否执行完毕
	}()
	wg.Wait()
}

// CloseBrowser 关闭浏览器
func (rodc *RODController) CloseBrowser() {
	rodc.Launcher.Kill()
	rodc.Browser = nil
	rodc.Launcher = nil
	log.Println("浏览器关闭")
}

/*以下方法拟废弃*/

// // GetCookiePath 获取平台cookie路径
// func (rodc *RODController) GetCookiePath(platformName string) (cookiePath string, err error) {
// 	cutl := utils.NewCommonUtils()

// 	// 获取默认配置文件夹
// 	configDir, err := cutl.GetConfigDir()
// 	if err != nil {
// 		return cookiePath, fmt.Errorf("获取配置文件夹错误: %w", err)
// 	}
// 	// 默认cookie路径
// 	cookiePath = filepath.Join(configDir, "cookies", platformName+".json")

// 	// 创建文件夹
// 	err = os.Mkdir(filepath.Dir(cookiePath), 0755)

// 	// 当不是文件存在错误的时候则抛出错误
// 	if err != nil && !os.IsExist(err) {
// 		return cookiePath, fmt.Errorf("创建cookies文件夹错误: %w", err)
// 	}

// 	return cookiePath, nil
// }

// // GetPlatformConfigPath 获取平台平台路径
// func (rodc *RODController) GetPlatformConfigPath(platformName string) (platformConfigPath string, err error) {

// 	cutl := utils.NewCommonUtils()

// 	// 获取默认配置文件夹
// 	configDir, err := cutl.GetConfigDir()
// 	if err != nil {
// 		return platformConfigPath, fmt.Errorf("获取配置文件夹错误: %w", err)
// 	}
// 	// 默认cookie路径
// 	platformConfigPath = filepath.Join(configDir, "platformConfig", platformName+".json")

// 	// 创建文件夹
// 	err = os.Mkdir(filepath.Dir(platformConfigPath), 0755)

// 	// 当不是文件存在错误的时候则抛出错误
// 	if err != nil && !os.IsExist(err) {
// 		return platformConfigPath, fmt.Errorf("创建platformConfig文件夹错误: %w", err)
// 	}

// 	return platformConfigPath, nil
// }

// // SaveCookies cookies保存为json
// func (rodc *RODController) SaveCookies(cookies []*proto.NetworkCookie, savePath string) (err error) {
// 	// 先序列化
// 	jsonStr, err := json.Marshal(cookies)
// 	if err != nil {
// 		return fmt.Errorf("JSON序列化失败: %w", err)
// 	}

// 	file, err := os.Create(savePath)
// 	if err != nil {
// 		return fmt.Errorf("创建文件夹失败: %w", err)
// 	}

// 	bw := bufio.NewWriter(file)
// 	_, err = bw.WriteString(string(jsonStr))
// 	if err != nil {
// 		return fmt.Errorf("保存JSON文件失败: %w", err)
// 	}
// 	bw.Flush()
// 	fmt.Println("Cookie 已经保存到本地")

// 	return nil

// }
