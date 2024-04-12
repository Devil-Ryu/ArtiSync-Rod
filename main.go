package main

import (
	"ArtiSync-Rod/backend/application"
	"ArtiSync-Rod/backend/controller"
	"ArtiSync-Rod/backend/platforms"
	"ArtiSync-Rod/backend/utils"
	"context"
	"embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	var err error

	/***创建控制器***/
	dbc := controller.NewDBController()  // 创建数据库连接
	rdc := controller.NewRODController() // 创建ROD控制器

	/***创建平台实例**/
	csdn := platforms.NewRodCSDN()   // 平台CSDN
	zhihu := platforms.NewRodZhiHu() // 平台CSDN

	/***创建App实例**/
	app := NewApp()

	/***创建文章App实例**/
	atApp := application.NewATApp()                // 文章APP
	atApp.SetController(dbc, rdc)                  // 设置文章APP控制器
	atApp.SetPlatforms([]interface{}{csdn, zhihu}) // 设置文章APP平台
	err = atApp.InitConfig()                       // 初始化文章APP配置
	if err != nil {
		log.Fatal(err)
	}

	/***创建工具类**/
	cutl := utils.NewCommonUtils() // 工具类

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "ArtiSync-Rod",
		Width:     1024,
		Height:    768,
		MinWidth:  1024,
		MinHeight: 768,
		// MaxWidth:          1280,
		// MaxHeight:         800,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         runtime.GOOS != "darwin",
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    nil,
			Middleware: nil,
		},
		Menu:     nil,
		Logger:   nil,
		LogLevel: logger.WARNING,
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			atApp.StartUp(ctx)
			cutl.SetContext(ctx)
		},
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			app,
			atApp,
			rdc,
			dbc,
			cutl,
			csdn,
			zhihu,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			DisableWindowIcon:    true,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
			BackdropType:        windows.Acrylic,
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHidden(),
			Appearance:           mac.NSAppearanceNameVibrantLight,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "ArtiSync-Rod",
				Message: "",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
