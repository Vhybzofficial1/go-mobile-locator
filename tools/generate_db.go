package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mobile-locator/internal/model"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 输入 JSON 网络地址
	inputURL := "https://raw.githubusercontent.com/zxc7563598/php-mobile-locator/main/src/data.json"
	// 输出 SQLite
	output := "internal/embedfiles/assets/carrier.db"
	// 删除旧文件
	_ = os.Remove(output)
	db, err := gorm.Open(sqlite.Open(output), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// 建表
	if err := db.AutoMigrate(&model.CarrierData{}); err != nil {
		log.Fatal(err)
	}
	// 请求网络 JSON
	resp, err := http.Get(inputURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("请求失败: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	t, err := decoder.Token()
	if err != nil || t != json.Delim('{') {
		log.Fatal("invalid JSON")
	}
	// 事务 + 批量
	tx := db.Begin()
	batch := make([]model.CarrierData, 0, 2000)
	total := 0
	for decoder.More() {
		keyTok, err := decoder.Token()
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		key := keyTok.(string)
		var info struct {
			Province string `json:"province"`
			City     string `json:"city"`
			ISP      string `json:"isp"`
		}
		if err := decoder.Decode(&info); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		batch = append(batch, model.CarrierData{
			Key:      key,
			Province: info.Province,
			City:     info.City,
			ISP:      info.ISP,
		})
		if len(batch) >= 2000 {
			if err := tx.CreateInBatches(batch, 2000).Error; err != nil {
				tx.Rollback()
				log.Fatal(err)
			}
			total += len(batch)
			fmt.Printf("插入了 %d 行...\n", total)
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		if err := tx.CreateInBatches(batch, 2000).Error; err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		total += len(batch)
	}
	tx.Commit()
	fmt.Printf("完成! 总计插入: %d\n", total)
}
