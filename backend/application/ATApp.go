package application

import (
	"ArtiSync-Rod/backend/controller"
	"ArtiSync-Rod/backend/db"
	"ArtiSync-Rod/backend/platforms"
	"ArtiSync-Rod/backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ATApp App
type ATApp struct {
	Ctx           context.Context
	Platforms     []interface{}             // 平台
	DBController  *controller.DBController  // 数据库控制器
	RODController *controller.RODController // 机器人控制器
	ArticleList   []utils.Article
}

// NewATApp 构造方法
func NewATApp() *ATApp {
	return &ATApp{}
}

// StartUp 设置context
func (at *ATApp) StartUp(ctx context.Context) {
	at.Ctx = ctx
}

// RefreshPlatforms 刷新平台
func (at *ATApp) RefreshPlatforms() (err error) {
	if at.HasController() == false {
		return fmt.Errorf("控制器未设置")
	}
	// 遍历平台，实例化平台控制器
	for _, param := range at.Platforms {
		switch platform := param.(type) {
		case *platforms.RodCSDN:
			log.Println("初始化CSDN平台")
			platform.InitRod(at.DBController, at.RODController, &platform.Config)
		case *platforms.RodZhiHu:
			log.Println("初始化知乎平台")
			platform.InitRod(at.DBController, at.RODController, &platform.Config)
		default:
			return fmt.Errorf("未知平台类型: %v", platform)
		}
	}
	return nil
}

// InitConfig 初始化配置
func (at *ATApp) InitConfig() (err error) {

	if at.HasController() == false {
		return fmt.Errorf("控制器未设置")
	}

	// 获取配置目录
	cutl := utils.NewCommonUtils()
	configDir, err := cutl.GetConfigDir()
	if err != nil {
		return err
	}

	// 连接数据库
	dbPath := path.Join(configDir, "artisync.db")
	err = at.DBController.Connect(dbPath)
	if err != nil {
		return err
	}

	// 刷新平台
	err = at.RefreshPlatforms()
	if err != nil {
		return err
	}

	return nil
}

// SetController 设置控制器
func (at *ATApp) SetController(dbController *controller.DBController, rodController *controller.RODController) {
	at.DBController = dbController
	at.RODController = rodController
}

// SetPlatforms 设置平台
func (at *ATApp) SetPlatforms(platforms []interface{}) {
	at.Platforms = platforms
}

// HasController 是否设置控制器
func (at *ATApp) HasController() bool {
	if at.DBController == nil || at.RODController == nil {
		return false
	}
	return true
}

// LoadArticles 加载文章
func (at *ATApp) LoadArticles(filePath string, imagePath string) (articleList []utils.Article, err error) {

	// 加载前清除数据
	at.ArticleList = []utils.Article{}

	files, err := os.ReadDir(filePath)
	if err != nil {
		err := fmt.Errorf("读取文章文件夹失败: %w", err)
		return at.ArticleList, err
	}

	// 遍历文件夹，获取文章信息
	for _, file := range files {
		if path.Ext(file.Name()) == ".md" {
			article := utils.Article{
				Title:          file.Name()[:len(file.Name())-3],
				Status:         utils.Waiting,
				SelectAccounts: []utils.AccountInfo{},
				// TODO(修复选择账号类型的问题)
			}

			// 设置文章路径以及图片路径以及图片读取方式
			article.MarkdownTool.MarkdownPath = path.Join(filePath, file.Name())
			article.MarkdownTool.ImagePath = imagePath

			// 开始分析文章
			err = article.MarkdownTool.AnalyzeMarkdown()
			if err != nil {
				err = fmt.Errorf("文章分析失败[%s]: %w", article.Title, err)
				return at.ArticleList, err
			}
			at.ArticleList = append(at.ArticleList, article)
		}
	}
	return at.ArticleList, nil
}

// SyncSelectAccounts 同步选择的平台
func (at *ATApp) SyncSelectAccounts(articles []utils.Article) {
	for index, article := range articles {
		at.ArticleList[index].SelectAccounts = []utils.AccountInfo{}
		// 遍历文章选择的平台，给平台信息增加初始信息，平台索引和选择平台列表中的索引一一对应
		for _, account := range article.SelectAccounts {
			at.ArticleList[index].SelectAccounts = append(at.ArticleList[index].SelectAccounts, utils.AccountInfo{
				ID: account.ID, Username: account.Username, PlatformKey: account.PlatformKey, PlatformAlias: account.PlatformAlias, Status: utils.PublishWating, Progress: 0, StepCount: 0,
			})
		}
		articles[index] = article
	}
}

// SyncPlatformsInfoFromRemote 同步平台信息(必要-统一方法)
func (at *ATApp) SyncPlatformsInfoFromRemote(remoteURL string) (err error) {
	// 下载链接文件
	content, err := utils.DownloadFile(remoteURL)
	if err != nil {
		return err
	}

	// 解析链接文件
	var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	var platforms []db.Platform
	err = mapstructure.Decode(data, &platforms)
	if err != nil {
		return fmt.Errorf("配置解码失败: %w", err)
	}

	// 保存平台信息到数据库
	if at.HasController() == false { // 如果没有数据库控制器则返回错误
		return fmt.Errorf("数据库控制器未设置")
	}
	_, err = at.DBController.CreatePlatforms(platforms)
	if err != nil {
		return err
	}

	return nil

}

// SyncSettingsFromRemote 同步设置(必要-统一方法)
func (at *ATApp) SyncSettingsFromRemote(remoteURL string) (err error) {
	// 下载链接文件
	content, err := utils.DownloadFile(remoteURL)
	if err != nil {
		return err
	}

	// 解析链接文件
	var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	var settings []db.Setting
	err = mapstructure.Decode(data, &settings)
	if err != nil {
		return fmt.Errorf("配置解码失败: %w", err)
	}

	// 保存设置信息到数据库
	if at.HasController() == false { // 如果没有数据库控制器则返回错误
		return fmt.Errorf("数据库控制器未设置")
	}
	_, err = at.DBController.CreateSettings(settings)
	if err != nil {
		return err
	}

	return nil

}

// GetArticlesInfo 获取文章信息
func (at *ATApp) GetArticlesInfo() []utils.Article {
	return at.ArticleList
}

// // genPlatformsInfo 生成平台信息
// func (at *ATApp) genPlatformsInfo() {
// 	// 遍历每个文章选择的平台，给平台信息增加初始信息，平台索引和选择平台列表中的索引一一对应
// 	for index := range at.ArticleList {
// 		for _, account := range at.ArticleList[index].SelectAccounts {
// 			at.ArticleList[index].AccountsInfo = append(at.ArticleList[index].AccountsInfo, utils.AccountInfo{
// 				ID: account.ID, PlatformKey: account.PlatformKey, PlatformAlias: account.PlatformAlias, Status: utils.PublishWating, Progress: 0, StepCount: 0,
// 			})
// 		}
// 	}
// }

func (at *ATApp) publishToAccount(account db.Account, article *utils.Article, accountProgress *utils.AccountInfo) (err error) {
	// 获取平台Key
	platormKey := account.PlatformKey

	// 选择对应平台
	switch platormKey {
	case "CSDN":
		platform := platforms.NewRodCSDN()
		err = platform.Start(at.Ctx, at.DBController, at.RODController, &platform.Config, &account, article, accountProgress, platform.Publish) // 实例化机器人
		if err != nil {
			return err
		}
	case "ZhiHu":
		platform := platforms.NewRodZhiHu()
		err = platform.Start(at.Ctx, at.DBController, at.RODController, &platform.Config, &account, article, accountProgress, platform.Publish) // 实例化机器人
		if err != nil {
			return err
		}
	}

	return nil
}

// Publish 发布文章
func (at *ATApp) Publish() (err error) {
	// at.genPlatformsInfo() // 生成要上传的平台的基础信息

	// 获取启用的账号信息
	// accountcs, err := at.DBController.QueryAccounts(map[string]interface{}{"Disabled": false})
	// log.Println("账号数量: ", len(accountcs))
	// if err != nil {
	// 	return err
	// }
	// 遍历文章
	for index := range at.ArticleList {
		fmt.Println("文章: ", at.ArticleList[index].Title)
		// 遍历账号
		publishCount := 0
		for accountIndex, accountProgress := range at.ArticleList[index].SelectAccounts {
			account, err := at.DBController.GetAccount(db.Account{ID: accountProgress.ID})
			if err != nil {
				return err
			}
			fmt.Println("账号: ", account)
			err = at.publishToAccount(account, &at.ArticleList[index], &at.ArticleList[index].SelectAccounts[accountIndex])
			if err != nil {
				return err
			}
			publishCount++
		}

		// 判断是否每个账号都发布完成
		if publishCount == len(at.ArticleList[index].SelectAccounts) {
			at.ArticleList[index].Status = utils.PublishedSuccess
			fmt.Println("文章发布成功: ", at.ArticleList[index].Title)
		} else {
			at.ArticleList[index].Status = utils.PublishedFailed
			fmt.Println("文章发布失败: ", at.ArticleList[index].Title)
		}

		runtime.EventsEmit(at.Ctx, "UpdatePlatformInfo")

	}

	return nil

}
