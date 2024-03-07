package structures

type Follow struct {
	ID             uint `json:"id" gorm:"primaryKey"`
	FollowerID     uint `json:"follower_id"`
	FollowedUserID uint `json:"followed_user_id"`
}
