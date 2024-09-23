package models

type Follower struct {
	ID         uint `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID uint `gorm:"not null" json:"follower_id"` // followしている人
	FolloweeID uint `gorm:"not null" json:"followee_id"` // followされている人

	// relations
	// Follower, Followee情報をFollowerデータと一緒に取得したい場合はPerloadを使用する
	Follower *User `gorm:"foreignKey:FollowerID;references:ID" json:"follower"`
	Followee *User `gorm:"foreignKey:FolloweeID;references:ID" json:"followee"`
}
