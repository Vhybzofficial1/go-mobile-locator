package dto

// Response 统一响应格式
type Response[T any] struct {
	Code    int    `json:"code"`    // 0 表示成功，非 0 表示失败
	Message string `json:"message"` // 成功或失败的提示
	Data    T      `json:"data"`    // 成功返回的数据，失败返回 nil
}

// Ok 返回成功结果
func Ok[T any](data T) Response[T] {
	return Response[T]{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
}

// Error 返回失败结果
func Error[T any](msg string) Response[T] {
	var empty T
	return Response[T]{
		Code:    -1,
		Message: msg,
		Data:    empty,
	}
}
