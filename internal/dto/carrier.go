package dto

// 手机号段数据通用结构
type CarrierData struct {
	Key      string `json:"key"`      // 号段
	Province string `json:"province"` // 省份
	City     string `json:"city"`     // 城市
	ISP      string `json:"isp"`      // 运营商
}

// 手机号段分页数据结构体
type CarrierPageData[T any] struct {
	Total int64         `json:"total"` // 总条数
	List  []CarrierData `json:"list"`  // 当前页数据
}

// 添加手机号段信息结构体
type CarrierCreateReq struct {
	Key      string `json:"key" binding:"required"`      // 号段
	Province string `json:"province" binding:"required"` // 省份
	City     string `json:"city" binding:"required"`     // 城市
	ISP      string `json:"isp" binding:"required"`      // 运营商
}

// 修改手机号段信息结构体
type CarrierUpdateReq struct {
	Province string `json:"province"` // 省份
	City     string `json:"city"`     // 城市
	ISP      string `json:"isp"`      // 运营商
}
