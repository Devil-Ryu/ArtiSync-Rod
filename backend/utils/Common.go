package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	sysruntime "runtime"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ResponseJSON 响应
type ResponseJSON struct {
	StatusCode int         // 状态码
	Message    string      // 返回信息
	Data       interface{} // 返回数据
}

// 状态枚举
const (
	Publishing       string = "发布中"  // 发布中（文章）
	PublishWating    string = "待发布"  // 待发布（文章）
	PublishedSuccess string = "发布成功" // 发布成功 （文章）
	PublishedFailed  string = "发布失败" // 发布失败（文章）
	RunningSuccess   string = "运行成功" // 运行成功（接口）
	RunningFailed    string = "运行失败" // 运行失败（接口）
	Waiting          string = "等待中"  // 等待中
	Running          string = "运行中"  // 运行中
)

// CommonUtils 常用工具
type CommonUtils struct {
	Ctx context.Context
}

// NewCommonUtils NewCommonUtils
func NewCommonUtils() *CommonUtils {
	return &CommonUtils{}

}

// SetContext 设置context
func (c *CommonUtils) SetContext(ctx context.Context) {
	c.Ctx = ctx
}

// OpenDir 打开文件夹
func (c *CommonUtils) OpenDir(defaultPath string) (selected string, err error) {
	opts := runtime.OpenDialogOptions{
		DefaultDirectory: filepath.Dir(defaultPath),
	}
	selected, err = runtime.OpenDirectoryDialog(c.Ctx, opts)

	// a.Print(artlog.DEBUG, "文章控制器", "选择文件夹: "+selected)
	return selected, err
}

// OpenFile 打开文件
func (c *CommonUtils) OpenFile(defaultPath string) (selected string, err error) {

	opts := runtime.OpenDialogOptions{
		DefaultDirectory: filepath.Dir(defaultPath),
	}
	selected, err = runtime.OpenFileDialog(c.Ctx, opts)

	return selected, err
}

// GetConfigDir 获取用户配置目录
func (c *CommonUtils) GetConfigDir() (dir string, err error) {
	systemType := sysruntime.GOOS
	if systemType == "darwin" {

	}
	// 获取用户配置目录
	dir, err = os.UserConfigDir()
	if err != nil {
		return dir, fmt.Errorf("获取用户配置目录失败: %w", err)
	}
	configDir := filepath.Join(dir, "ArtiSync-Rod")

	err = os.Mkdir(configDir, 0755)

	// 当不是文件存在错误的时候则抛出错误
	// 忽略文件夹存在错误
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		return "", err
	}

	return configDir, err
}

// GetConfigFilePath 获取配置文件
func (c *CommonUtils) GetConfigFilePath() (configFilePath string, err error) {
	// 获取配置文件地址
	configPath, err := c.GetConfigDir()
	if err != nil {
		return configFilePath, err
	}
	configFilePath = filepath.Join(configPath, "config.json")

	defaultConfig := map[string]interface{}{
		"imagePath":       "",
		"imageSelect":     "相对文章目录",
		"selectPlatforms": []string{},
	}

	// 尝试加载配置文件
	_, err = c.LoadJSONFile(configFilePath)
	if err != nil {
		// 将默认配置进行保存
		err = c.SaveJSONFile(configFilePath, defaultConfig)
	}
	return configFilePath, err
}

// LoadJSONFile 加载JSON文件
func (c *CommonUtils) LoadJSONFile(fileDir string) (jsonMap map[string]interface{}, err error) {
	dataBytes, err := os.ReadFile(fileDir)
	if err != nil {
		return jsonMap, fmt.Errorf("读取JSON文件失败: %w", err)
	}

	// 读取文件内容
	err = json.Unmarshal(dataBytes, &jsonMap)
	if err != nil {
		return jsonMap, fmt.Errorf("JSON文件序列化失败: %w", err)
	}
	return jsonMap, nil
}

// SaveJSONFile 加载JSON文件
func (c *CommonUtils) SaveJSONFile(fileDir string, jsonMap map[string]interface{}) (err error) {
	// 先序列化

	jsonStr, err := json.Marshal(jsonMap)
	if err != nil {
		return fmt.Errorf("保存JSON文件失败: %w", err)
	}

	// 格式化字符串
	var formatStr bytes.Buffer
	err = json.Indent(&formatStr, jsonStr, "", "\t")
	if err != nil {
		return fmt.Errorf("JSON文件格式化失败: %w", err)
	}

	file, err := os.Create(fileDir)
	if err != nil {
		return fmt.Errorf("创建文件夹失败: %w", err)
	}
	bw := bufio.NewWriter(file)
	_, err = bw.WriteString(formatStr.String())
	if err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}
	bw.Flush()
	return nil
}

// ReadJSONFile 函数用于读取并解析指定路径的JSON文件
func (c *CommonUtils) ReadJSONFile(filePath string) (interface{}, error) {
	// 尝试读取文件内容
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法读取文件: %v", err)
	}

	// 检测文件内容是否为空
	if len(fileContent) == 0 {
		return nil, fmt.Errorf("文件内容为空")
	}

	// 定义用于存储解析后数据的结构体
	var data interface{}

	// 尝试解析JSON数据
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, fmt.Errorf("无法解析JSON数据: %v", err)
	}

	return data, nil
}

// DownloadFile 下载图片
func DownloadFile(downloadURL string) ([]byte, error) {
	// 请求链接
	response, err := http.Get(downloadURL)
	if err != nil {
		err = fmt.Errorf("网络请求错误: %w", err)
		return nil, err
	}
	defer response.Body.Close()

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("读取网络响应错误: %w", err)
		return nil, err
	}

	// 关闭连接

	return respBytes, nil
}
