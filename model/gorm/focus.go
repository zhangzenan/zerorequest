package gorm

// 定义轮播图的结构体
type Focus struct {
	Id       string
	Title    string
	Pic      string
	Link     string
	Position int
	Status   int
}

// 配置操作数据库的表
func (Focus) TableName() string {
	return "focus"
}
