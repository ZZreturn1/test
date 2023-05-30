package database

// 导入所需要的包
import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/fs"
	"os"
	"path"
	"x-ui/config"
	"x-ui/database/model"
)

// 变量 db 为gorm.DB类型对象的指针
var db *gorm.DB

// initUser 初始化用户表
func initUser() error {
                // 自动迁移用户表
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

                // 查询用户表记录数
	var count int64
	err = db.Model(&model.User{}).Count(&count).Error
	if err != nil {
		return err
	}

                // 如果用户表记录数为0，则创建默认管理员用户
	if count == 0 {
		user := &model.User{
			Username: "admin",
			Password: "admin",
		}
		return db.Create(user).Error
	}

                // 返回相应的零值或错误信息
	return nil
}

// initInbound 初始化入站表
func initInbound() error {
                // 调用 db 对象的 AutoMigrate 方法，该方法用于执行数据库迁移操作
	return db.AutoMigrate(&model.Inbound{})
}

// initSetting 初始化设置表
func initSetting() error {
                // 执行数据库迁移操作，将 model.Setting 结构体对应的表格创建或更新到数据库中
	return db.AutoMigrate(&model.Setting{})
}

// InitDB 初始化数据库
func InitDB(dbPath string) error {
                // 创建数据库文件所在的目录
	dir := path.Dir(dbPath)	// 变量dir
	err := os.MkdirAll(dir, fs.ModeDir)	// 变量err
	if err != nil {
		return err
	}

	var gormLogger logger.Interface

                // 根据是否处于调试模式选择日志记录方式
	if config.IsDebug() {
		gormLogger = logger.Default
	} else {
		gormLogger = logger.Discard
	}

                // 配置 GORM
	c := &gorm.Config{
		Logger: gormLogger,
	}

                // 打开数据库连接
	db, err = gorm.Open(sqlite.Open(dbPath), c)
	if err != nil {
		return err
	}

                // 初始化用户表
	err = initUser()
	if err != nil {
		return err
	}

                // 初始化入站表
	err = initInbound()
	if err != nil {
		return err
	}

                // 初始化设置表
	err = initSetting()
	if err != nil {
		return err
	}

                // 返回相应的零值或错误信息
	return nil
}

// GetDB 返回数据库连接对象
func GetDB() *gorm.DB {
	return db
}

// IsNotFound 检查错误是否为记录未找到的错误
func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
