package test

import (
	"ArtiSync-Rod/backend/platforms"
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	// app := application.NewATApp()
	// app.LoadArticles("/Users/ryu/Documents/test", "/Users/ryu/Documents/test")
	zhihu := platforms.NewRodZhiHu()
	// err := zhihu.Login()
	// fmt.Println("Login err: ", err)

	autInfo, err := zhihu.CheckAuthentication()
	if err != nil {
		log.Fatal("err2:", err)
	} else {
		log.Println("autInfo:", autInfo)
	}
	// csdn.SetArticle(&app.ArticleList[0])
	// csdn.OpenCSDNPage("https://editor.csdn.net/md/?articleId=137137778")

	// fmt.Println("csdn.Article.Progress:", csdn.Article.Progress)
	// fmt.Println("app.ArticleList[0].Progress:", app.ArticleList[0].Progress)
	// app.SyncSelectPlatforms(app.ArticleList)
	// for index := range app.ArticleList {
	// 	fmt.Println("index", index)
	// }

}
