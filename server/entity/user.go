package entity

type User struct {
	ID       uint64   `gorm:"primary_key:auto_increment" json:"id"`
	Username string   `gorm:"uniqueIndex;type:varchar(255)" json:"username"`
	Password string   `gorm:"->;<-;not null" json:"-"`
	Token    string   `gorm:"-" json:"token,omitempty"`
	Posts    *[]Posts `json:"posts,omitempty"`
}
