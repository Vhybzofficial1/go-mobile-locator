package service

// Service 层的依赖注入容器
//
// 用于集中管理各个 Service 实例，方便在 Handler 或 App 层统一访问。
// 通过 Container，可以轻松替换或 Mock 各个 Service。
type Container struct {
	Carrier *CarrierService
}
