package main

import (
	"embed"
	"log"
	"mobile-locator/internal/config"
	"mobile-locator/internal/repository"
	"mobile-locator/internal/service"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化配置
	if err := config.InitConfig(""); err != nil {
		log.Fatal("初始化配置失败:", err)
	}
	// 初始化数据库
	if err := config.InitDB(); err != nil {
		log.Fatal("初始化数据库失败:", err)
	}
	defer config.CloseDB()
	// 依赖注入
	db := config.GetDB()
	carrierRepo := repository.NewCarrierRepository(db)
	svcContainer := &service.Container{
		Carrier: service.NewCarrierService(carrierRepo),
	}
	// 创建应用程序
	app := NewApp(svcContainer)
	err := wails.Run(&options.App{
		Title:  "手机号归属地查询",
		Width:  1280,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
