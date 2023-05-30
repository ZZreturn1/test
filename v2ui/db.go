package v2ui

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var v2db *gorm.DB

func initDB(dbPath string) error {
	c := &gorm.Config{
		// 配置日志记录器为丢弃模式，即不输出日志
		Logger: logger.Discard, 
	}
	var err error
	// 打开 SQLite 数据库连接
	v2db, err = gorm.Open(sqlite.Open(dbPath), c) 
	if err != nil {
		// 如果打开数据库连接失败，则返回错误
		return err 
	}

	return nil
}

func getV2Inbounds() ([]*V2Inbound, error) {
	// 创建 V2Inbound 切片
	inbounds := make([]*V2Inbound, 0) 

	// 查询数据库中的 V2Inbound 记录并将结果存入 inbounds 切片
	err := v2db.Model(V2Inbound{}).Find(&inbounds).Error 

	// 返回查询结果和错误信息
	return inbounds, err 
}
