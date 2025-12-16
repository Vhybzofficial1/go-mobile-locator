package main

import (
	"context"
	"mobile-locator/internal/dto"
	"mobile-locator/internal/service"
	"os"

	wailsappRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 声明应用的核心结构体
// 相当于整个桌面应用的“入口对象”，
// 前端调用的 Go 方法，最终都会挂在这个结构体上
type App struct {
	// Wails 在应用启动时传入的上下文
	// 用于调用 runtime 相关方法（如窗口控制、事件、对话框等）
	ctx context.Context
	// 服务容器，内部组合了各个业务 Service
	// App 本身不处理具体业务逻辑，而是把请求转交给 Service 层
	svc *service.Container
}

// 在应用启动时自动调用的方法
// 该方法只会执行一次
// 这里保存启动时传入的 context，
// 以便在后续方法中使用 Wails runtime 提供的能力
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// 用于创建 App 实例
// 在这里通过参数的方式注入 Service 容器，
// 属于应用启动阶段的“依赖注入”
func NewApp(svc *service.Container) *App {
	return &App{svc: svc}
}

// CarrierCreate 添加手机号段信息
// 供前端调用，通过 Service 层完成数据创建
func (a *App) CarrierCreate(req dto.CarrierCreateReq) dto.Response[any] {
	if err := a.svc.Carrier.Create(a.ctx, req); err != nil {
		return dto.Error[any](err.Error())
	}
	return dto.Ok[any](nil)
}

// CarrierGet 根据手机号段查询单条数据
// 供前端调用，返回对应的运营商信息
func (a *App) CarrierGet(key string) dto.Response[*dto.CarrierData] {
	data, err := a.svc.Carrier.Get(a.ctx, key[:7])
	if err != nil {
		return dto.Error[*dto.CarrierData](err.Error())
	}
	return dto.Ok(data)
}

// CarrierUpdate 根据手机号段查询并修改单条数据
// 供前端调用，通过 Service 层对数据进行修改
func (a *App) CarrierUpdate(key string, req dto.CarrierUpdateReq) dto.Response[any] {
	if err := a.svc.Carrier.Update(a.ctx, key, req); err != nil {
		return dto.Error[any](err.Error())
	}
	return dto.Ok[any](nil)
}

// CarrierDelete 根据手机号段查询并删除单条数据
// 供前端调用，通过 Service 层对数据进行删除
func (a *App) CarrierDelete(key string) dto.Response[any] {
	if err := a.svc.Carrier.Delete(a.ctx, key); err != nil {
		return dto.Error[any](err.Error())
	}
	return dto.Ok[any](nil)
}

// CarrierList 分页查询手机号段列表
// 供前端调用，通过 Service 层模糊查询所有数据并进行分野
// 如果 key 存在，则会筛选与 key 相关的信息
func (a *App) CarrierList(key string, page int, size int) dto.Response[dto.CarrierPageData[dto.CarrierData]] {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	list, total, err := a.svc.Carrier.List(a.ctx, key, page, size)
	if err != nil {
		return dto.Error[dto.CarrierPageData[dto.CarrierData]](err.Error())
	}
	return dto.Ok(dto.CarrierPageData[dto.CarrierData]{
		Total: total,
		List:  list,
	})
}

// CarrierProcessCSV 批量解析 CSV 文件并查询手机号段信息
// 供前端调用，接收一个 Base64 编码的 CSV 文件，CSV 第一列为手机号
// 方法通过 Service 层查询数据库中对应手机号段信息，并填充省份、城市、运营商等字段
// 返回处理后的 CSV 二进制数据，由前端用于下载
func (a *App) CarrierProcessCSV(base64Str string) dto.Response[[]byte] {
	data, err := a.svc.Carrier.ProcessCSV(a.ctx, base64Str)
	if err != nil {
		return dto.Error[[]byte](err.Error())
	}
	return dto.Ok(data)
}

// SaveCSV 保存处理后的 CSV 文件到本地
// 供前端调用，参数 data 为 CSV 文件的二进制数据，通常来自 CarrierProcessCSV 方法处理后的结果
// 方法会弹出保存文件对话框，允许用户选择保存路径和文件名
// 如果用户取消保存（path 为空），方法直接返回 nil
// 文件保存失败会返回对应的错误
func (a *App) SaveCSV(data []byte) error {
	path, err := wailsappRuntime.SaveFileDialog(a.ctx, wailsappRuntime.SaveDialogOptions{
		Title:           "保存处理后的 CSV 文件",
		DefaultFilename: "result.csv",
		Filters: []wailsappRuntime.FileFilter{
			{DisplayName: "CSV 文件", Pattern: "*.csv"},
		},
	})
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	return os.WriteFile(path, data, 0644)
}
