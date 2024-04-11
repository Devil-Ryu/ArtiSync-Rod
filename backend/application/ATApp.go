package application

import (
	"ArtiSync-Rod/backend/controller"
	"ArtiSync-Rod/backend/db"
	"ArtiSync-Rod/backend/platforms"
	"ArtiSync-Rod/backend/utils"
	"context"
	"fmt"
	"log"
	"os"
	"path"
)

// ATApp App
type ATApp struct {
	Ctx           context.Context
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

	return nil
}

// SetController 设置控制器
func (at *ATApp) SetController(dbController *controller.DBController, rodController *controller.RODController) {
	at.DBController = dbController
	at.RODController = rodController
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
				Title:           file.Name()[:len(file.Name())-3],
				Status:          utils.Waiting,
				SelectPlatforms: []string{},
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

// SyncSelectPlatforms 同步选择的平台
func (at *ATApp) SyncSelectPlatforms(data []utils.Article) {
	for i := 0; i < len(data); i++ {
		at.ArticleList[i].SelectPlatforms = data[i].SelectPlatforms
	}
}

// GetArticlesInfo 获取文章信息
func (at *ATApp) GetArticlesInfo() []utils.Article {
	return at.ArticleList
}

// genPlatformsInfo 生成平台信息
func (at *ATApp) genPlatformsInfo() {
	// 遍历每个文章选择的平台，给平台信息增加初始信息，平台索引和选择平台列表中的索引一一对应
	for index := range at.ArticleList {
		for _, platformName := range at.ArticleList[index].SelectPlatforms {
			at.ArticleList[index].PlatformsInfo = append(at.ArticleList[index].PlatformsInfo, utils.PlatformInfo{
				Name: platformName, Status: utils.PublishWating, Progress: 0, StepCount: 0,
			})
		}
	}
}

func (at *ATApp) publishToAccount(account db.Account, article *utils.Article) (err error) {
	// 获取平台Key
	platormKey := account.PlatformKey

	// 选择对应平台
	switch platormKey {
	case "CSDN":
		platform := platforms.NewRodCSDN()
		err = platform.Start(at.DBController, at.RODController, &platform.Config, account, article, platform.Publish) // 实例化机器人
		if err != nil {
			return err
		}
	case "ZhiHu":
		platform := platforms.NewRodZhiHu()
		err = platform.Start(at.DBController, at.RODController, &platform.Config, account, article, platform.Publish) // 实例化机器人
		if err != nil {
			return err
		}
	}

	return nil
}

// Publish 发布文章
func (at *ATApp) Publish() (err error) {
	at.genPlatformsInfo() // 生成要上传的平台的基础信息

	// 获取启用的账号信息
	accountcs, err := at.DBController.QueryAccounts(map[string]interface{}{"Disabled": false})
	log.Println("账号数量: ", len(accountcs))
	if err != nil {
		return err
	}

	// 遍历文章
	for index := range at.ArticleList {
		fmt.Println("文章: ", at.ArticleList[index].Title)
		// 遍历账号
		for _, account := range accountcs {
			fmt.Println("账号: ", account)
			err = at.publishToAccount(account, &at.ArticleList[index])
			if err != nil {
				return err
			}
		}

	}

	return nil

}
