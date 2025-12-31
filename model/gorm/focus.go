package gorm

// 定义轮播图的结构体
type Focus struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Pic      string `json:"pic"`
	Link     string `json:"link"`
	Position int    `json:"position"`
	Status   int    `json:"status"`
}

// 配置操作数据库的表
func (Focus) TableName() string {
	return "focus"
}
