package engines

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MysqlConf struct {
	// mysql相关配置
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	ShowSql      bool   `yaml:"showSql"`
	MysqlMaxOpen int    `yaml:"mysql_max_open"`
	MysqlMaxIdle int    `yaml:"mysql_max_idle"`
}

func CreateMysqlEngine(dbInfo MysqlConf, plugins ...gorm.Plugin) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(getMysqlConnURL(dbInfo)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if dbInfo.MysqlMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(dbInfo.MysqlMaxIdle)
	}
	if dbInfo.MysqlMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(dbInfo.MysqlMaxOpen)
	}
	for _, plugin := range plugins {
		err = db.Use(plugin)
		if err != nil {
			return nil, err
		}
	}

	db.Config.Logger = db.Logger.LogMode(logger.Info)

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 获取数据库连接的url
func getMysqlConnURL(info MysqlConf) (url string) {
	url = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		info.User,
		info.Password,
		info.Host,
		info.Port,
		info.Database,
		info.Charset)
	return
}
