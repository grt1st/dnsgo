package backends

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DSN = "linker:root123456@tcp(127.0.0.1:3306)/mars?charset=utf8mb4&parseTime=True&loc=Local"
	DB  *gorm.DB
)

func InitDB() error {
	var err error
	DB, err = gorm.Open(mysql.Open(DSN))
	if err != nil {
		panic(err)
	}
	return DB.AutoMigrate(DNSConfig{}, DNSRecord{}, HTTPRecord{})
}

type DNSConfig struct {
	gorm.Model
	Name  string
	Value string
	Kind  string
}

func (DNSConfig) TableName() string {
	return "dns_config"
}

type DNSRecord struct {
	gorm.Model
	ClientIP string
	Name     string
	Value    string
	Kind     string
}

func (DNSRecord) TableName() string {
	return "dns_record"
}

type HTTPRecord struct {
	gorm.Model
	ClientIP   string
	ReqHeaders string
	ReqBody    string
	URL        string
}

func (HTTPRecord) TableName() string {
	return "http_record"
}

type Setting struct {
	gorm.Model
	Name string
	Value string
}

func (Setting) TableName() string {
	return "kv_settings"
}
