package models

type User struct {
	ID       uint    `gorm:"primaryKey" json:"usersID"`
	Username string  `json:"username" binding:"required"`
	Email    string  `json:"email" binding:"required;unique"`
	Password string  `json:"password" binding:"required,min=6"`
	Photos   []Photo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
