package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CarrierData 用户模型
type CarrierData struct {
	Key       string         `gorm:"primaryKey;size:7;comment:7位数字key, 例如:1300001"`
	Province  string         `gorm:"size:20;not null;index;comment:省份"`
	City      string         `gorm:"size:20;not null;index;comment:城市"`
	ISP       string         `gorm:"size:10;not null;comment:运营商"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// TableName 表名
func (CarrierData) TableName() string {
	return "carrier_data"
}

// Validate 验证数据有效性
func (m *CarrierData) Validate() error {
	if len(m.Key) != 7 {
		return fmt.Errorf("key必须是7位数字")
	}
	for _, c := range m.Key {
		if c < '0' || c > '9' {
			return fmt.Errorf("key只能包含数字")
		}
	}
	if m.Province == "" {
		return fmt.Errorf("省份不能为空")
	}
	if m.City == "" {
		return fmt.Errorf("城市不能为空")
	}
	if m.ISP == "" {
		return fmt.Errorf("运营商不能为空")
	}
	return nil
}
