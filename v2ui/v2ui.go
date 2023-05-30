package v2ui

import (
	"fmt"
	"x-ui/config"
	"x-ui/database"
	"x-ui/database/model"
	"x-ui/util/common"
	"x-ui/web/service"
)

func MigrateFromV2UI(dbPath string) error {
	err := initDB(dbPath)
	if err != nil {
		return common.NewError("初始化 v2-ui 数据库失败:", err)
	}
	err = database.InitDB(config.GetDBPath())
	if err != nil {
		return common.NewError("初始化 x-ui 数据库失败:", err)
	}

	v2Inbounds, err := getV2Inbounds()
	if err != nil {
		return common.NewError("获取 v2-ui 入站配置失败:", err)
	}
	if len(v2Inbounds) == 0 {
		fmt.Println("迁移 v2-ui 入站配置成功：0")
		return nil
	}

	userService := service.UserService{}
	user, err := userService.GetFirstUser()
	if err != nil {
		return common.NewError("获取 x-ui 用户失败:", err)
	}

	inbounds := make([]*model.Inbound, 0)
	for _, v2inbound := range v2Inbounds {
		inbounds = append(inbounds, v2inbound.ToInbound(user.Id))
	}

	inboundService := service.InboundService{}
	err = inboundService.AddInbounds(inbounds)
	if err != nil {
		return common.NewError("添加 x-ui 入站配置失败:", err)
	}

	fmt.Println("迁移 v2-ui 入站配置成功:", len(inbounds))

	return nil
}

