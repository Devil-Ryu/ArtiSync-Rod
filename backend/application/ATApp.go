package application

import (
	"ArtiSync-Rod/backend/platforms"
	"ArtiSync-Rod/backend/utils"
	"context"
	"fmt"
	"os"
	"path"
)

// ATApp App
type ATApp struct {
	Ctx         context.Context
	ArticleList []utils.Article
}

// NewATApp 构造方法
func NewATApp() *ATApp {
	return &ATApp{}
}

// StartUp 设置context
func (at *ATApp) StartUp(ctx context.Context) {
	at.Ctx = ctx
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

func (at *ATApp) runPlatform(platformName string, platformIndex int, article *utils.Article) (err error) {
	switch platformName {
	case "CSDN":
		bot := platforms.NewRodCSDN()
		bot.Init(at.Ctx, article, platformIndex)
		err = bot.RUN()
		if err != nil {
			return err
		}
		fmt.Println("CSDN Done")
	case "ZhiHu":
		bot := platforms.NewRodZhiHu()
		bot.Init(at.Ctx, article, platformIndex)
		err = bot.RUN()
		if err != nil {
			return err
		}

	}

	return nil
}

// Run Run
func (at *ATApp) Run() (err error) {
	at.genPlatformsInfo() // 生成要上传的平台的基础信息

	// 遍历文章
	for index := range at.ArticleList {
		fmt.Println("文章: ", at.ArticleList[index].Title)
		// 遍历平台
		for platformIndex, platforName := range at.ArticleList[index].SelectPlatforms {

			fmt.Println("平台: ", platforName)
			err = at.runPlatform(platforName, platformIndex, &at.ArticleList[index])
			if err != nil {
				return err
			}
		}
	}

	return nil

}
