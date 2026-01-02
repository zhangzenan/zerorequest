package gorm

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  int    `json:"sex"`
}

func (User) TableName() string {
	return "user"
}
