package models

// User representa un usuario en el sistema
type User struct {
	ID          int    `gorm:"primaryKey;column:id" json:"id"`
	FirebaseUID string `gorm:"uniqueIndex;not null" json:"firebase_uid"`
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	IsAdmin     bool   `gorm:"default:false" json:"isAdmin"`
	CreatedAt   string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   string `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "appa_users"
}
