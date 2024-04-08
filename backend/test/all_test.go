package test

import (
	"ArtiSync-Rod/backend/application"
	"ArtiSync-Rod/backend/platforms"
	"log"
	"testing"
)

func TestAll(t *testing.T) {

	// err := zhihu.Login()
	// fmt.Println("Login err: ", err)

	// csdn.SetArticle(&app.ArticleList[0])
	// csdn.OpenCSDNPage("https://editor.csdn.net/md/?articleId=137137778")

	// fmt.Println("csdn.Article.Progress:", csdn.Article.Progress)
	// fmt.Println("app.ArticleList[0].Progress:", app.ArticleList[0].Progress)
	// app.SyncSelectPlatforms(app.ArticleList)
	// for index := range app.ArticleList {
	// 	fmt.Println("index", index)
	// }

}

func TestAuth(t *testing.T) {
	// 测试认证接口
	zhihu := platforms.NewRodZhiHu()
	autInfo, err := zhihu.CheckAuthentication()
	if err != nil {
		log.Fatal("err2:", err)
	} else {
		log.Println("autInfo:", autInfo)
	}
}

func TestRun(t *testing.T) {
	app := application.NewATApp()
	app.LoadArticles("/Users/ryu/Documents/test", "/Users/ryu/Documents/test")
	zhihu := platforms.NewRodZhiHu()
	zhihu.Init(app.Ctx, &app.ArticleList[1], 0)
	zhihu.RUN()
}

func TestCheckConfig(t *testing.T) {
	zhihu := platforms.NewRodZhiHu()
	err := zhihu.SetConfig(false)
	err2 := zhihu.SetConfig(false)
	log.Println("err: ", err)
	log.Println("err2: ", err2)
}

func TestOpenPage(t *testing.T) {
	zhihu := platforms.NewRodZhiHu()
	zhihu.SetConfig(false)
	zhihu.OpenPage(zhihu.Config.ArticleManagePage)
}
