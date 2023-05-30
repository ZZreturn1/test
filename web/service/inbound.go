package service

import (
	"fmt"
	"time"
	"x-ui/database"
	"x-ui/database/model"
	"x-ui/util/common"
	"x-ui/xray"

	"gorm.io/gorm"
)

// InboundService是入站服务结构体
type InboundService struct {
}

// GetInbounds获取指定用户的入站信息列表
func (s *InboundService) GetInbounds(userId int) ([]*model.Inbound, error) {
	db := database.GetDB()	// 获取数据库连接
	var inbounds []*model.Inbound	// 入站信息列表
	err := db.Model(model.Inbound{}).Where("user_id = ?", userId).Find(&inbounds).Error	/ 查询指定用户的入站信息
	// 处理查询错误
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return inbounds, nil	// 返回入站信息列表
}

// GetAllInbounds获取所有入站信息列表
func (s *InboundService) GetAllInbounds() ([]*model.Inbound, error) {
	db := database.GetDB()	// 获取数据库连接
	var inbounds []*model.Inbound	// 入站信息列表
	err := db.Model(model.Inbound{}).Find(&inbounds).Error	// 查询所有入站信息
	// 处理查询错误
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return inbounds, nil	// 返回入站信息列表
}

// checkPortExist检查指定端口是否已存在
func (s *InboundService) checkPortExist(port int, ignoreId int) (bool, error) {
	db := database.GetDB()	// 获取数据库连接
	db = db.Model(model.Inbound{}).Where("port = ?", port)		// 根据端口查询入站信息
	// 如果ignoreId大于0，则忽略指定ID的记录
	if ignoreId > 0 {
		db = db.Where("id != ?", ignoreId)
	}

	var count int64
	err := db.Count(&count).Error	// 统计查询结果数量
	// 处理查询错误
	if err != nil {
		return false, err
	}
	return count > 0, nil	// 返回端口是否已存在的结果
}

// AddInbound添加入站信息
func (s *InboundService) AddInbound(inbound *model.Inbound) error {
	exist, err := s.checkPortExist(inbound.Port, 0)	// 检查端口是否已存在
	// 处理检查错误
	if err != nil {
		return err
	}
	// 如果端口已存在，则返回错误
	if exist {
		return common.NewError("端口已存在:", inbound.Port)
	}
	db := database.GetDB()	// 获取数据库连接
	return db.Save(inbound).Error	// 保存入站信息
}

// AddInbounds批量添加入站信息
func (s *InboundService) AddInbounds(inbounds []*model.Inbound) error {
	// 遍历入站信息列表
	for _, inbound := range inbounds {
		exist, err := s.checkPortExist(inbound.Port, 0)	// 检查端口是否已存在
		// 处理检查错误
		if err != nil {
			return err
		}
		// 如果端口已存在，则返回错误
		if exist {
			return common.NewError("端口已存在:", inbound.Port)
		}
	}

	db := database.GetDB()	// 获取数据库连接
	tx := db.Begin()	// 开启事务
	var err error
	// 延迟执行事务提交或回滚操作
	defer func() {
		if err == nil {
			tx.Commit()	// 提交事务
		} else {
			tx.Rollback()	// 回滚事务
		}
	}()
	
	// 遍历入站信息列表
	for _, inbound := range inbounds {
		err = tx.Save(inbound).Error	// 保存入站信息
		if err != nil {
			return err		// 处理保存错误
		}
	}

	return nil		// 返回nil表示操作成功
}

// DelInbound删除指定ID的入站信息
func (s *InboundService) DelInbound(id int) error {
	db := database.GetDB()	// 获取数据库连接
	return db.Delete(model.Inbound{}, id).Error	// 删除指定ID的入站信息
}

// GetInbound获取指定ID的入站信息
func (s *InboundService) GetInbound(id int) (*model.Inbound, error) {
	db := database.GetDB()	// 获取数据库连接
	inbound := &model.Inbound{}	// 入站信息对象
	err := db.Model(model.Inbound{}).First(inbound, id).Error		// 查询指定ID的入站信息
	// 处理查询错误
	if err != nil {
		return nil, err
	}
	return inbound, nil		// 返回入站信息对象
}

// UpdateInbound更新入站信息
func (s *InboundService) UpdateInbound(inbound *model.Inbound) error {
	// 检查端口是否已存在
	exist, err := s.checkPortExist(inbound.Port, inbound.Id)
	// 处理检查错误
	if err != nil {
		return err
	}
	// 如果端口已存在，则返回错误
	if exist {
		return common.NewError("端口已存在:", inbound.Port)
	}

	oldInbound, err := s.GetInbound(inbound.Id)	// 获取旧的入站信息
	// 处理查询错误 
	if err != nil {
		return err
	}
	// 更新入站信息的各个字段
	oldInbound.Up = inbound.Up
	oldInbound.Down = inbound.Down
	oldInbound.Total = inbound.Total
	oldInbound.Remark = inbound.Remark
	oldInbound.Enable = inbound.Enable
	oldInbound.ExpiryTime = inbound.ExpiryTime
	oldInbound.Listen = inbound.Listen
	oldInbound.Port = inbound.Port
	oldInbound.Protocol = inbound.Protocol
	oldInbound.Settings = inbound.Settings
	oldInbound.StreamSettings = inbound.StreamSettings
	oldInbound.Sniffing = inbound.Sniffing
	oldInbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)

	db := database.GetDB()	// 获取数据库连接
	return db.Save(oldInbound).Error	// 更新入站信息
}

// AddTraffic添加流量统计
func (s *InboundService) AddTraffic(traffics []*xray.Traffic) (err error) {
	// 如果流量列表为空，则直接返回nil
	if len(traffics) == 0 {
		return nil
	}
	db := database.GetDB()	// 获取数据库连接
	db = db.Model(model.Inbound{})	// 操作入站信息表
	tx := db.Begin()	// 开启事务
	// 延迟执行事务提交或回滚操作
	defer func() {
		if err != nil {
			tx.Rollback()	// 回滚事务
		} else {
			tx.Commit()	// 提交事务
		}
	}()
	// 遍历流量列表
	for _, traffic := range traffics {
		// 如果流量是入站流量
		if traffic.IsInbound {
			err = tx.Where("tag = ?", traffic.Tag).	// 根据标签查询入站信息
				UpdateColumn("up", gorm.Expr("up + ?", traffic.Up)).	// 更新入站流量的上行数据量
				UpdateColumn("down", gorm.Expr("down + ?", traffic.Down)).	// 更新入站流量的下行数据量
				Error	// 处理更新错误
			// 如果更新出现错误，则直接返回
			if err != nil {
				return
			}
		}
	}
	return
}

// DisableInvalidInbounds禁用无效的入站信息
func (s *InboundService) DisableInvalidInbounds() (int64, error) {
	db := database.GetDB()	// 获取数据库连接
	now := time.Now().Unix() * 1000	// 当前时间戳（毫秒）
	result := db.Model(model.Inbound{}).	// 操作入站信息表
		Where("((total > 0 and up + down >= total) or (expiry_time > 0 and expiry_time <= ?)) and enable = ?", now, true).		 // 根据条件查询无效的入站信息
		Update("enable", false)	// 将无效的入站信息的enable字段设置为false
	err := result.Error	// 处理更新错误
	count := result.RowsAffected	// 获取更新的记录数量
	return count, err	// 返回更新的记录数量和错误信息
}
