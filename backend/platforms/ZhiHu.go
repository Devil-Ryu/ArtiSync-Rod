package platforms

import (
	"ArtiSync-Rod/backend/controller"
	"ArtiSync-Rod/backend/db"
	"ArtiSync-Rod/backend/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// RodZhiHu ZhiHu 机器人
type RodZhiHu struct {
	*Model
	Config *ConfigZhiHu // 配置
}

// ConfigZhiHu 配置文件
type ConfigZhiHu struct {
	Disabled                      bool   `json:"Disabled" `                                         // 是否禁用(有默认值的不用设置validate，否则0值会报错)
	LoginPageURL                  string `json:"LoginPageURL" validate:"required"`                  // 登录页
	HomePageURL                   string `json:"HomePageURL" validate:"required"`                   // 首页
	CreateArticlePage             string `json:"CreateArticlePage" validate:"required"`             // 创建文章页
	ArticleManagePage             string `json:"ArticleManagePage" validate:"required"`             // 文章管理页
	ProfilePageURL                string `json:"ProfilePageURL" validate:"required"`                // 个人页
	LoginBoxSelector              string `json:"LoginBoxSelector" validate:"required"`              // 登录页登录box选择器
	TitleInputSelector            string `json:"TitleInputSelector" validate:"required"`            // 标题输入选择器
	ContentAreaSelector           string `json:"ContentAreaSelector" validate:"required"`           // 内容区域选择器
	ImageUploadStep1Selector      string `json:"ImageUploadStep1Selector" validate:"required"`      // 图片上传选择器
	ImageUploadStep2Selector      string `json:"ImageUploadStep2Selector" validate:"required"`      // 图片上传选择器
	ImageUploadStep3Selector      string `json:"ImageUploadStep3Selector" validate:"required"`      // 图片上传选择器
	UploadArticleBtnStep1Selector string `json:"UploadArticleBtnStep1Selector" validate:"required"` // 保存文章1选择器
	UploadArticleBtnStep2Selector string `json:"UploadArticleBtnStep2Selector" validate:"required"` // 保存文章2选择器
	UploadArticleBtnStep3Selector string `json:"UploadArticleBtnStep3Selector" validate:"required"` // 保存文章3选择器
	ProfileNameSelector           string `json:"ProfileNameSelector" validate:"required"`           // 个人简介名称选择器
	ProfileAvatarSelector         string `json:"ProfileAvatarSelector" validate:"required"`         // 个人简介头像选择器
}

// NewRodZhiHu 初始化
func NewRodZhiHu() *RodZhiHu {
	return &RodZhiHu{Model: &Model{Key: "ZhiHu", Alias: "知乎"}}
}

// Login 登录CSDN后把cookie保存到本地（重写方法）
func (zhihu *RodZhiHu) Login() (err error) {
	// 检查基础配置
	err = zhihu.CheckConfig(zhihu.Config)
	if err != nil {
		return err
	}

	// 打开一个新的浏览器用作登录
	rdc := controller.NewRODController()
	rdc.StartBrowser(false) // 显示浏览器
	// 确认浏览器关闭
	defer rdc.CloseBrowser()

	var loginCookies []*proto.NetworkCookie
	// 访问登录页面
	rdc.Browser.MustPage(zhihu.Config.LoginPageURL)
	fmt.Println("zhihu.Config.LoginPageURL", zhihu.Config.LoginPageURL)

	// 监听是否登录成功
	for {
		time.Sleep(time.Duration(rdc.CheckTime) * time.Second) // 监听频率
		pages, _ := rdc.Browser.Pages()
		log.Println("当前页面数: ", len(pages), " , 登录状态: 等待登录")
		// 如果页面全部关闭，则推出
		if len(pages) == 0 {
			log.Println("当前页面数: ", len(pages), " , 登录状态: 取消登录")
			break
		} else {
			targetPage, err := pages.FindByURL(zhihu.Config.HomePageURL)
			if err == nil && targetPage.MustInfo().URL == zhihu.Config.HomePageURL {
				log.Println("当前页面数: ", len(pages), " , 登录状态: 登录成功")
				loginCookies = targetPage.MustCookies() // 获取cookie，终止监听

				cookieParams := proto.CookiesToParams(loginCookies)

				// 存入数据库
				zhihu.DBController.CreateOrUpdateAccounts([]db.Account{{
					PlatformKey:   zhihu.Key,
					PlatformAlias: zhihu.Alias,
					Username:      "", // 手动登录的默认没有
					LoginType:     "", // 手动登录的默认没有
					Password:      "", // 手动登录的默认没有
					Cookies:       cookieParams}})
				break
			} else {

			}
		}
	}

	return err
}

// CheckAuthentication 检查是否授权（重写方法）
func (zhihu *RodZhiHu) CheckAuthentication() (authInfo map[string]string, err error) {
	// 检查基础配置
	err = zhihu.CheckConfig(zhihu.Config)
	if err != nil {
		return authInfo, err
	}

	// 确认是否有账号
	if zhihu.HasAccount() == false {
		return authInfo, fmt.Errorf(zhihu.Alias + "账号未设置")
	}

	/*设置浏览器*/
	if zhihu.RODController.Browser == nil {
		zhihu.RODController.StartBrowser(false) // 启动浏览器，无头模式
	}

	// 确认浏览器关闭
	defer zhihu.RODController.CloseBrowser()

	profileURL := zhihu.Config.ProfilePageURL

	zhihu.RODController.Browser.SetCookies(zhihu.Account.Cookies)
	page := zhihu.RODController.Browser.MustPage()

	// 导航到页面，判断是否超时
	err = rod.Try(func() {
		page.Timeout(6 * time.Second).MustNavigate(profileURL)

	})

	if errors.Is(err, context.DeadlineExceeded) {
		return authInfo, fmt.Errorf("超时错误: %w", err)
	} else if err != nil {
		return authInfo, fmt.Errorf("其他错误: %w", err)
	}

	// 轮询查询到哪个页面
	page.Race().Element(zhihu.Config.LoginBoxSelector).MustHandle(func(e *rod.Element) {
		err = fmt.Errorf("Cookie失效")
		return

	}).Element(zhihu.Config.ProfileNameSelector).MustHandle(func(e *rod.Element) {
		name, err := page.MustWaitStable().MustElement(zhihu.Config.ProfileNameSelector).Text()
		if err != nil {
			err = fmt.Errorf("获取Name失败: %w", err)
			return
		}
		authInfo = map[string]string{
			"ID":   "",
			"name": name,
		}
	}).MustDo()

	log.Println(authInfo)
	return authInfo, nil

}

// Publish 发布文章（重写方法）
func (zhihu *RodZhiHu) Publish() (err error) {
	log.Println("开始运行: " + zhihu.Alias)

	// 检查基础配置
	err = zhihu.CheckConfig(zhihu.Config)
	if err != nil {
		zhihu.Article.Status = utils.PublishedFailed
		return err
	}

	// 检查是否有文章
	if zhihu.Article == nil {
		return fmt.Errorf("文章未设置")
	}

	// 确认是否有账号
	if zhihu.HasAccount() == false {
		zhihu.Article.Status = utils.PublishedFailed
		return fmt.Errorf(zhihu.Alias + "账号未设置")
	}
	zhihu.Article.Status = utils.Publishing

	/*设置浏览器*/
	if zhihu.RODController.Browser == nil {
		zhihu.RODController.StartBrowser(false) // 启动浏览器
	}
	// 确认浏览器关闭
	defer zhihu.RODController.CloseBrowser()

	/*设置浏览器cookies*/
	err = zhihu.RODController.Browser.SetCookies(zhihu.Account.Cookies)
	if err != nil {
		zhihu.Article.Status = utils.PublishedFailed
		// runtime.EventsEmit(zhihu.Ctx, "UpdatePlatformInfo")
		return err
	}

	/*创建文章*/
	page := zhihu.RODController.Browser.MustPage(zhihu.Config.CreateArticlePage) // 导航到写文章页

	/*上传图片*/
	for index, imageInfo := range zhihu.Article.MarkdownTool.ImagesInfo { // 遍历图片列表，上传图片
		imagePath := path.Join(zhihu.Article.MarkdownTool.ImagePath, imageInfo.URL)
		uploadURL, err := zhihu.uploadImage(page, imagePath)
		if err != nil {
			zhihu.Article.Status = utils.PublishedFailed
			// runtime.EventsEmit(zhihu.Ctx, "UpdatePlatformInfo")
			return err
		}
		zhihu.Article.MarkdownTool.ImagesInfo[index].UploadURL = uploadURL
		zhihu.UpdatePlatformInfo()
	}
	/*替换图片*/
	zhihu.Article.MarkdownTool.ReplaceImages()
	/*保存替换图片后的文章到本地*/
	savePath, err := zhihu.Article.MarkdownTool.SaveToMarkdown()
	if err != nil {
		zhihu.Article.Status = utils.PublishedFailed
		// runtime.EventsEmit(zhihu.Ctx, "UpdatePlatformInfo")
		return fmt.Errorf("保存Markdown失败: %w", err)
	}
	/*上传文章*/
	page.MustElement(zhihu.Config.UploadArticleBtnStep1Selector).MustClick()                             // 点击导入文章按钮
	page.MustElement(zhihu.Config.UploadArticleBtnStep2Selector).MustClick()                             // 点击导入文章按钮
	page.MustWaitStable().MustElement(zhihu.Config.UploadArticleBtnStep3Selector).MustSetFiles(savePath) // 导入本地文章

	/*设置文章标题*/
	page.MustElement(zhihu.Config.TitleInputSelector).MustSelectAllText().MustInput(zhihu.Article.Title)

	zhihu.UpdatePlatformInfo()

	zhihu.Article.Status = utils.PublishedSuccess
	// 获取URL并更新
	// zhihu.Article.PlatformsInfo[zhihu.PlatformIndex].PublishURL = zhihu.Config.ArticleManagePage
	// runtime.EventsEmit(zhihu.Ctx, "UpdatePlatformInfo")
	page.MustWaitStable()

	return nil

}

// UpdatePlatformInfo 更新平台上传进度(默认进度包括图片上传，文件上传的步骤)(TODO 解决Progress为int的问题0)
func (zhihu *RodZhiHu) UpdatePlatformInfo() {

	// 更新文章中平台上传的进度
	// zhihu.Article.PlatformsInfo[zhihu.PlatformIndex].StepCount++
	// zhihu.Article.PlatformsInfo[zhihu.PlatformIndex].Progress = float32(zhihu.Article.PlatformsInfo[ZhiHu.PlatformIndex].StepCount) / float32(len(ZhiHu.Article.MarkdownTool.ImagesInfo)+1) * 100 // +1是因为后面还有一个上传文章

	// // 更新文章总上传进度
	zhihu.Article.Progress = 0
	for _, platformInfo := range zhihu.Article.PlatformsInfo {
		zhihu.Article.Progress += platformInfo.Progress
	}
	zhihu.Article.Progress = zhihu.Article.Progress / float32(len(zhihu.Article.PlatformsInfo))

	// runtime.EventsEmit(zhihu.Ctx, "UpdatePlatformInfo")
}

/****************************自定义函数区****************************/
// clearContent 清除内容
func (zhihu *RodZhiHu) clearContent(element *rod.Element) {
	// 找到所有的输入
	subEls := element.MustElementsX("//div[contains(@class, 'css-5sjb75')]")
	log.Println("subEls: ", len(subEls))

	// 清空输入区
	for _, subEl := range subEls {
		subEl.MustClick()
	}
}

// UploadImage 上传图片
func (zhihu *RodZhiHu) uploadImage(page *rod.Page, imagePath string) (uploadURL string, err error) {
	contentAreaEl := page.MustElement(zhihu.Config.ContentAreaSelector) // 定位内容区域
	// 访问上传图片页面
	page.MustElement(zhihu.Config.ImageUploadStep1Selector).MustClick()
	page.MustWaitStable().MustElement(zhihu.Config.ImageUploadStep2Selector).MustSetFiles(imagePath)
	page.MustWaitStable().MustElement(zhihu.Config.ImageUploadStep3Selector).MustClick()

	// 获取图片元素
	imgEl, err := contentAreaEl.MustWaitStable().ElementX("*//img")

	if err != nil {
		return uploadURL, err
	}
	// 获取图片链接

	uploadURL = *imgEl.MustAttribute("src")
	fmt.Println("uploadUrl: ", uploadURL)
	// 获得链接后清除内容
	zhihu.clearContent(contentAreaEl)
	return uploadURL, nil
}
