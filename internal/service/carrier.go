package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"io"
	"runtime"
	"sync"
	"time"

	"mobile-locator/internal/dto"
	"mobile-locator/internal/model"
	"mobile-locator/internal/repository"
)

type CarrierService struct {
	repo repository.CarrierRepository
}

func NewCarrierService(repo repository.CarrierRepository) *CarrierService {
	return &CarrierService{repo: repo}
}

// Create 添加一个新的手机号段信息
//
// 将前端请求的手机号段信息插入数据库。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - req: CarrierCreateReq DTO，包含手机号段、所属省份、城市和运营商
//
// 返回值：
//   - error: 创建过程中可能产生的错误
func (s *CarrierService) Create(ctx context.Context, req dto.CarrierCreateReq) error {
	return s.repo.Create(ctx, &model.CarrierData{
		Key:      req.Key,
		Province: req.Province,
		City:     req.City,
		ISP:      req.ISP,
	})
}

// Get 根据手机号段获取单条信息
//
// 查询数据库中指定手机号段的记录，并转换为 DTO 返回。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - key: 手机号段（前 7 位）
//
// 返回值：
//   - *dto.CarrierData: 查询到的手机号段信息
//   - error: 查询过程中可能产生的错误或不存在记录时的错误
func (s *CarrierService) Get(ctx context.Context, key string) (*dto.CarrierData, error) {
	m, err := s.repo.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("该号段无数据")
	}
	return &dto.CarrierData{
		Key:      m.Key,
		Province: m.Province,
		City:     m.City,
		ISP:      m.ISP,
	}, nil
}

// Update 更新指定手机号段的信息
//
// 仅更新请求中非空的字段，同时更新 updated_at 时间戳。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - key: 手机号段（前 7 位）
//   - req: CarrierUpdateReq DTO，包含要更新的省份、城市和运营商
//
// 返回值：
//   - error: 更新过程中可能产生的错误
func (s *CarrierService) Update(ctx context.Context, key string, req dto.CarrierUpdateReq) error {
	updates := map[string]any{
		"updated_at": time.Now(),
	}
	if req.Province != "" {
		updates["province"] = req.Province
	}
	if req.City != "" {
		updates["city"] = req.City
	}
	if req.ISP != "" {
		updates["isp"] = req.ISP
	}
	return s.repo.Update(ctx, key, updates)
}

// Delete 删除指定手机号段信息
//
// 从数据库中删除指定手机号段的记录。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - key: 手机号段（前 7 位）
//
// 返回值：
//   - error: 删除过程中可能产生的错误
func (s *CarrierService) Delete(ctx context.Context, key string) error {
	return s.repo.Delete(ctx, key)
}

// List 分页查询手机号段信息
//
// 根据手机号段关键字进行模糊查询，并返回分页数据。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - key: 手机号段关键字，用于模糊查询
//   - page: 页码，从 1 开始
//   - size: 每页记录数
//
// 返回值：
//   - []dto.CarrierData: 当前页的手机号段记录列表
//   - int64: 总记录数
//   - error: 查询过程中可能产生的错误
func (s *CarrierService) List(ctx context.Context, key string, page, size int) ([]dto.CarrierData, int64, error) {
	list, total, err := s.repo.List(ctx, key, page, size)
	if err != nil {
		return nil, 0, err
	}
	res := make([]dto.CarrierData, len(list))
	for i, m := range list {
		res[i] = dto.CarrierData{
			Key:      m.Key,
			Province: m.Province,
			City:     m.City,
			ISP:      m.ISP,
		}
	}
	return res, total, nil
}

// ListAll 获取所有手机号段信息
//
// 该方法从数据库中查询所有手机号段记录，并将其转换为 DTO 格式返回，方便在前端使用。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//
// 返回值：
//   - []dto.CarrierData: 包含全部手机号段信息的切片，每条记录包含手机号段、所属省份、城市和运营商
//   - error: 查询或转换过程中可能产生的错误
func (s *CarrierService) ListAll(ctx context.Context) ([]dto.CarrierData, error) {
	list, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]dto.CarrierData, len(list))
	for i, m := range list {
		res[i] = dto.CarrierData{
			Key:      m.Key,
			Province: m.Province,
			City:     m.City,
			ISP:      m.ISP,
		}
	}
	return res, nil
}

// ProcessCSV 解析上传的 CSV 文件（Base64 编码）并填充手机号段信息
//
// 该方法的功能流程如下：
// 1. 将前端上传的 Base64 编码 CSV 文件解码为二进制数据。
// 2. 读取 CSV 表头，保留前四列（手机号、省份、城市、运营商），其余列原样保留。
// 3. 从数据库中加载所有手机号段信息到内存 Map，用于快速查询。
// 4. 使用多线程（按 CPU 核心数）处理每一行数据：
//   - 取出手机号前 7 位作为 key 查询内存 Map，获取对应的省份、城市和运营商。
//   - 将查询结果填充到 CSV 行中，并保留原有其余列。
//
// 5. 保证输出的 CSV 行顺序与输入一致。
// 6. 将处理后的 CSV 写入内存缓冲，并以二进制形式返回。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - base64Str: 前端上传的 CSV 文件 Base64 编码字符串
//
// 返回值：
//   - []byte: 处理完成的 CSV 文件二进制数据，可直接返回给前端下载
//   - error: 处理过程中可能产生的错误
func (s *CarrierService) ProcessCSV(ctx context.Context, base64Str string) ([]byte, error) {
	// 解码 Base64
	fileData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bytes.NewReader(fileData))
	outBuf := &bytes.Buffer{}
	writer := csv.NewWriter(outBuf)
	// 读取表头
	headerRow, err := reader.Read()
	if err != nil {
		return nil, err
	}
	header := []string{"手机号", "省份", "城市", "运营商"}
	if len(headerRow) > 4 {
		header = append(header, headerRow[4:]...)
	}
	writer.Write(header)
	// 加载数据库到内存 Map
	dbList, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	cache := make(map[string]*model.CarrierData, len(dbList))
	for i := range dbList {
		item := dbList[i]
		cache[item.Key] = &item
	}
	// 多线程处理行
	type resultRow struct {
		index int
		row   []string
	}
	rowsCh := make(chan struct {
		index int
		row   []string
	}, 1000)
	resultsCh := make(chan resultRow, 1000)
	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range rowsCh {
				row := item.row
				if len(row) == 0 {
					resultsCh <- resultRow{item.index, row}
					continue
				}
				phone := row[0]
				restCols := []string{}
				if len(row) > 4 {
					restCols = row[4:]
				}
				var province, city, isp string
				if len(phone) >= 7 {
					key := phone[:7]
					if data, ok := cache[key]; ok && data != nil {
						province = data.Province
						city = data.City
						isp = data.ISP
					}
				}
				newRow := []string{phone, province, city, isp}
				newRow = append(newRow, restCols...)
				resultsCh <- resultRow{item.index, newRow}
			}
		}()
	}
	// 读取 CSV 行并发送给 workers
	go func() {
		index := 0
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}
			rowsCh <- struct {
				index int
				row   []string
			}{index, row}
			index++
		}
		close(rowsCh)
	}()
	// 收集结果并保证顺序输出
	go func() {
		wg.Wait()
		close(resultsCh)
	}()
	buffered := make(map[int][]string)
	nextIndex := 0
	for res := range resultsCh {
		buffered[res.index] = res.row
		for {
			if row, ok := buffered[nextIndex]; ok {
				writer.Write(row)
				delete(buffered, nextIndex)
				nextIndex++
			} else {
				break
			}
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}
