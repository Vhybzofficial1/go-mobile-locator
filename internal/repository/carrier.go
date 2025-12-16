package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"mobile-locator/internal/model"
)

type CarrierRepository interface {
	Create(ctx context.Context, m *model.CarrierData) error
	Get(ctx context.Context, key string) (*model.CarrierData, error)
	Update(ctx context.Context, key string, updates map[string]interface{}) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, key string, page, size int) ([]model.CarrierData, int64, error)
	ListAll(ctx context.Context) ([]model.CarrierData, error)
}

type carrierRepository struct {
	db *gorm.DB
}

func NewCarrierRepository(db *gorm.DB) CarrierRepository {
	return &carrierRepository{db: db}
}

// Create 在数据库中添加一条手机号段记录
//
// 如果记录已存在且被软删除，则会恢复并更新信息；
// 如果记录已存在且未删除，则返回错误。
// 否则创建新记录。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - m: CarrierData 模型，包含手机号段、所属省份、城市和运营商
//
// 返回值：
//   - error: 创建或恢复过程中可能产生的错误
func (r *carrierRepository) Create(ctx context.Context, m *model.CarrierData) error {
	db := r.db.WithContext(ctx)
	var exist model.CarrierData
	err := db.
		Unscoped().
		Where("key = ?", m.Key).
		First(&exist).Error
	if err == nil {
		if exist.DeletedAt.Valid {
			return db.Unscoped().
				Model(&exist).
				Updates(map[string]any{
					"province":   m.Province,
					"city":       m.City,
					"isp":        m.ISP,
					"deleted_at": nil,
				}).Error
		}
		return fmt.Errorf("该号码段已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return db.Create(m).Error
}

// Get 根据手机号段获取单条数据库记录
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - key: 手机号段（前 7 位）
//
// 返回值：
//   - *model.CarrierData: 查询到的记录，如果不存在返回 nil
//   - error: 查询过程中可能产生的错误
func (r *carrierRepository) Get(ctx context.Context, key string) (*model.CarrierData, error) {
	var m model.CarrierData
	err := r.db.WithContext(ctx).First(&m, "key = ?", key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

// Update 更新指定手机号段的数据库记录
//
// 根据 key 定位记录，并根据 updates 更新字段。
// updates 可以包含 province、city、isp 等字段。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - key: 手机号段（前 7 位）
//   - updates: 字段更新映射
//
// 返回值：
//   - error: 更新过程中可能产生的错误
func (r *carrierRepository) Update(ctx context.Context, key string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&model.CarrierData{}).
		Where("key = ?", key).
		Updates(updates).Error
}

// Delete 删除指定手机号段的数据库记录
//
// 从数据库中删除 key 对应的记录（软删除或物理删除取决于 GORM 配置）。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - key: 手机号段（前 7 位）
//
// 返回值：
//   - error: 删除过程中可能产生的错误
func (r *carrierRepository) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).
		Where("key = ?", key).
		Delete(&model.CarrierData{}).Error
}

// List 分页查询手机号段记录
//
// 根据 key 进行模糊查询（可匹配手机号段、省份、城市或运营商），
// 并返回分页数据。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//   - key: 查询关键字，可为空，支持模糊匹配
//   - page: 页码，从 1 开始
//   - size: 每页记录数
//
// 返回值：
//   - []model.CarrierData: 当前页查询结果列表
//   - int64: 总记录数
//   - error: 查询过程中可能产生的错误
func (r *carrierRepository) List(ctx context.Context, key string, page, size int) ([]model.CarrierData, int64, error) {
	var list []model.CarrierData
	var total int64
	query := r.db.WithContext(ctx).Model(&model.CarrierData{})
	if key != "" {
		likeKey := "%" + key + "%"
		query = query.Where(
			"key LIKE ? OR province LIKE ? OR city LIKE ? OR isp LIKE ?",
			likeKey, likeKey, likeKey, likeKey,
		)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	err := query.Order("key").Offset(offset).Limit(size).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ListAll 获取数据库中所有手机号段记录
//
// 返回所有记录列表，一般用于内存加载或批量处理。
// 注意：数据量大时可能会占用较多内存。
//
// 参数：
//   - ctx: 上下文，用于控制请求生命周期
//
// 返回值：
//   - []model.CarrierData: 所有手机号段记录
//   - error: 查询过程中可能产生的错误
func (r *carrierRepository) ListAll(ctx context.Context) ([]model.CarrierData, error) {
	var list []model.CarrierData
	err := r.db.WithContext(ctx).Find(&list).Error
	return list, err
}
