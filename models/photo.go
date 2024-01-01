package models

type Photo struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
	UserID   uint   `json:"user_id"`
}
