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

// RodCSDN CSDN机器人
type RodCSDN struct {
	*Model
	Config *ConfigCSDN // 配置
}

// ConfigCSDN 配置文件
type ConfigCSDN struct {
	Disabled                 bool   `json:"Disabled" `                                    // 是否禁用(有默认值的不用设置validate，否则0值会报错)
	LoginPageURL             string `json:"LoginPageURL" validate:"required"`             // 登录页
	HomePageURL              string `json:"HomePageURL" validate:"required"`              // 首页
	ArticleManagePage        string `json:"ArticleManagePage" validate:"required"`        // 文章管理页
	ProfilePageURL           string `json:"ProfilePageURL" validate:"required"`           // 个人页
	LoginBoxSelector         string `json:"LoginBoxSelector" validate:"required"`         // 登录页登录box选择器
	CreateArticleBtnSelector string `json:"CreateArticleBtnSelector" validate:"required"` // 创建文章选择器
	TitleInputSelector       string `json:"TitleInputSelector" validate:"required"`       // 标题输入选择器
	ContentAreaSelector      string `json:"ContentAreaSelector" validate:"required"`      // 内容区域选择器
	ImageUploadStep1Selector string `json:"ImageUploadStep1Selector" validate:"required"` // 图片上传选择器
	ImageUploadStep2Selector string `json:"ImageUploadStep2Selector" validate:"required"` // 图片上传选择器
	UploadArticleBtnSelector string `json:"UploadArticleBtnSelector" validate:"required"` // 保存文章选择器
	SaveArticleBtnSelector   string `json:"SaveArticleBtnSelector" validate:"required"`   // 保存文章选择器
	ProfileIDSelector        string `json:"ProfileIDSelector" validate:"required"`        // 个人简介ID选择器
	ProfileNameSelector      string `json:"ProfileNameSelector" validate:"required"`      // 个人简介名称选择器
	ProfileAvatarSelector    string `json:"ProfileAvatarSelector" validate:"required"`    // 个人简介头像选择器
}

// NewRodCSDN 初始化
func NewRodCSDN() *RodCSDN {
	return &RodCSDN{Model: &Model{Key: "CSDN", Alias: "CSDN"}}
}

// CheckAuthentication 检查是否授权（重写方法）
func (csdn *RodCSDN) CheckAuthentication() (authInfo map[string]string, err error) {
	// 检查基础配置
	err = csdn.CheckConfig(csdn.Config)
	if err != nil {
		return authInfo, err
	}

	// 确认是否有账号
	if csdn.HasAccount() == false {
		return authInfo, fmt.Errorf(csdn.Alias + "账号未设置")
	}

	/*设置浏览器*/
	if csdn.RODController.Browser == nil {
		csdn.RODController.StartBrowser(false) // 启动浏览器，无头模式
	}

	// 确认浏览器关闭
	defer csdn.RODController.CloseBrowser()

	profileURL := csdn.Config.ProfilePageURL

	csdn.RODController.Browser.SetCookies(csdn.Account.Cookies)
	page := csdn.RODController.Browser.MustPage()

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
	page.Race().Element(csdn.Config.LoginBoxSelector).MustHandle(func(e *rod.Element) {
		err = fmt.Errorf("Cookie失效")

	}).Element(csdn.Config.ProfileNameSelector).MustHandle(func(e *rod.Element) {
		name, err := page.MustWaitStable().MustElement(csdn.Config.ProfileNameSelector).Text()
		if err != nil {
			err = fmt.Errorf("获取Name失败: %w", err)
			return
		}
		ID, err := page.MustElement(csdn.Config.ProfileIDSelector).Text()
		if err != nil {
			err = fmt.Errorf("获取ID失败: %w", err)
			return
		}
		authInfo = map[string]string{
			"ID":   ID,
			"name": name,
		}
	}).MustDo()

	log.Println(authInfo)
	return authInfo, err
}

// Login 登录CSDN后把cookie保存到本地（重写方法）
func (csdn *RodCSDN) Login() (err error) {
	// 检查基础配置
	err = csdn.CheckConfig(csdn.Config)
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
	rdc.Browser.MustPage(csdn.Config.LoginPageURL)

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
			targetPage, err := pages.FindByURL(csdn.Config.HomePageURL)
			if err == nil {
				log.Println("当前页面数: ", len(pages), " , 登录状态: 登录成功")
				loginCookies = targetPage.MustCookies() // 获取cookie，终止监听
				cookieParams := proto.CookiesToParams(loginCookies)

				// 存入数据库
				csdn.DBController.CreateOrUpdateAccounts([]db.Account{{
					PlatformKey:   csdn.Key,
					PlatformAlias: csdn.Alias,
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

// Publish 发布文章（重写方法）
func (csdn *RodCSDN) Publish() (err error) {
	// 检查基础配置
	err = csdn.CheckConfig(csdn.Config)
	if err != nil {
		csdn.Article.Status = utils.PublishedFailed
		return err
	}

	// 检查是否有文章
	if csdn.Article == nil {
		return fmt.Errorf("文章未设置")
	}

	// 确认是否有账号
	if csdn.HasAccount() == false {
		csdn.Article.Status = utils.PublishedFailed
		return fmt.Errorf("CSDN账号未设置")
	}

	csdn.Article.Status = utils.Publishing

	/*设置浏览器*/
	if csdn.RODController.Browser == nil {
		csdn.RODController.StartBrowser(false) // 启动浏览器
	}
	// 确认浏览器关闭
	defer csdn.RODController.CloseBrowser()

	/*设置浏览器cookies*/
	err = csdn.RODController.Browser.SetCookies(csdn.Account.Cookies)
	if err != nil {
		csdn.Article.Status = utils.PublishedFailed
		// runtime.EventsEmit(csdn.Ctx, "UpdatePlatformInfo")
		return err
	}

	/*创建文章*/
	page := csdn.RODController.Browser.MustPage(csdn.Config.HomePageURL)                               // 导航到首页
	page.MustElement(csdn.Config.CreateArticleBtnSelector).MustClick()                                 // 光标移动到发布文章按钮
	contentAreaEl := page.MustElement(csdn.Config.ContentAreaSelector)                                 // 定位内容区域
	csdn.clearContent(contentAreaEl)                                                                   // 清除现有输入区内容
	page.MustElement(csdn.Config.TitleInputSelector).MustSelectAllText().MustInput(csdn.Article.Title) // 输入文章标题

	/*上传图片*/
	for index, imageInfo := range csdn.Article.MarkdownTool.ImagesInfo { // 遍历图片列表，上传图片
		imagePath := path.Join(csdn.Article.MarkdownTool.ImagePath, imageInfo.URL)
		uploadURL, err := csdn.uploadImage(page, imagePath)
		if err != nil {
			csdn.Article.Status = utils.PublishedFailed
			// runtime.EventsEmit(csdn.Ctx, "UpdatePlatformInfo")
			return err
		}
		csdn.Article.MarkdownTool.ImagesInfo[index].UploadURL = uploadURL
		csdn.UpdatePlatformInfo()
	}
	/*替换图片*/
	csdn.Article.MarkdownTool.ReplaceImages()
	savePath, err := csdn.Article.MarkdownTool.SaveToMarkdown() // 保存到本地
	if err != nil {
		csdn.Article.Status = utils.PublishedFailed
		// runtime.EventsEmit(csdn.Ctx, "UpdatePlatformInfo")
		return fmt.Errorf("保存Markdown失败: %w", err)
	}
	/*上传文章*/
	page.MustElement(csdn.Config.UploadArticleBtnSelector).MustSetFiles(savePath) // 导入本地文章

	csdn.UpdatePlatformInfo()
	page.MustElement(csdn.Config.SaveArticleBtnSelector).MustClick() // 点击保存

	csdn.Article.Status = utils.PublishedSuccess
	// 获取URL并更新
	// csdn.Article.PlatformsInfo[csdn.PlatformIndex].PublishURL = csdn.Config.ArticleManagePage
	// runtime.EventsEmit(csdn.Ctx, "UpdatePlatformInfo")

	return nil

}

// UpdatePlatformInfo 更新平台上传进度(默认进度包括图片上传，文件上传的步骤)(TODO 解决Progress为int的问题0)
func (csdn *RodCSDN) UpdatePlatformInfo() {

	// 更新文章中平台上传的进度
	// csdn.Article.PlatformsInfo[csdn.PlatformIndex].StepCount++
	// csdn.Article.PlatformsInfo[csdn.PlatformIndex].Progress = float32(csdn.Article.PlatformsInfo[csdn.PlatformIndex].StepCount) / float32(len(csdn.Article.MarkdownTool.ImagesInfo)+1) * 100 // +1是因为后面还有一个上传文章

	// 更新文章总上传进度
	csdn.Article.Progress = 0
	for _, platformInfo := range csdn.Article.PlatformsInfo {
		csdn.Article.Progress += platformInfo.Progress
	}
	csdn.Article.Progress = csdn.Article.Progress / float32(len(csdn.Article.PlatformsInfo))

	// runtime.EventsEmit(csdn.Ctx, "UpdatePlatformInfo")
}

/****************************自定义函数区****************************/
// clearContent 清除内容
func (csdn *RodCSDN) clearContent(element *rod.Element) {
	// 找到所有的输入
	subEls := element.MustElementsX(".//div")

	// 清空输入区
	for _, subEl := range subEls {
		subEl.Remove()
	}
}

// UploadImage 上传图片
func (csdn *RodCSDN) uploadImage(page *rod.Page, imagePath string) (uploadURL string, err error) {
	contentAreaEl := page.MustElement(csdn.Config.ContentAreaSelector) // 定位内容区域
	// 访问上传图片页面
	page.MustElement(csdn.Config.ImageUploadStep1Selector).MustClick()
	page.MustWaitStable().MustElement(csdn.Config.ImageUploadStep2Selector).MustSetFiles(imagePath).WaitInvisible()

	// 获取图片元素
	imgEl, err := contentAreaEl.MustWaitStable().ElementX("*//img")
	if err != nil {
		return uploadURL, err
	}
	// 获取图片链接
	uploadURL = *imgEl.MustAttribute("src")
	fmt.Println("uploadUrl: ", uploadURL)
	// 获得链接后清除内容
	csdn.clearContent(contentAreaEl)
	return uploadURL, nil
}
