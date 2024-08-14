package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255)" json:"password"`
	Dob       time.Time `gorm:"type:date" json:"dob"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
