package model

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func initTable(tables ...interface{}) {
	for _, table := range tables {
		// 如果数据库中不存在该表，则创建表
		if !db.Migrator().HasTable(table) {
			// 创建表时添加后缀
			if err := db.Set("gorm:table_options",
				"ENGINE=InnoDB DEFAULT CHARSET=utf8").
				Migrator().CreateTable(table); err != nil {
				zap.L().Error(err.Error())
			}
		}
	}
}

type PaginationQ struct {
	//每页显示的数量
	PageSize int `json:"page_size"`
	//当前页码
	Page int `json:"page"`
	//分页的数据内容
	Data interface{} `json:"data"`
	//全部的页码数量
	Total int64 `json:"total"`
}

// 分页扫描器
func (p *PaginationQ) PaginateScan(queryTx *gorm.DB) (data *PaginationQ, err error) {
	err = queryTx.Count(&p.Total).Error
	if err != nil {
		return p, err
	}
	switch {
	case p.PageSize > 100:
		p.PageSize = 100
	case p.PageSize < 1:
		p.PageSize = 1
	}
	offset := (p.Page - 1) * p.PageSize
	err = queryTx.Offset(offset).Limit(p.PageSize).Scan(p.Data).Error
	return p, err
}
