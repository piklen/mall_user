package serializer

import (
	"time"
	"user/model"
)

type User struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	NickName string `json:"nickname"`
	Type     int    `json:"type"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
	CreateAt int64  `json:"create_at"`
}

// BuildUser 序列化用户
func BuildUser(user *model.User) map[string]interface{} {
	return map[string]interface{}{
		"ID":       user.ID,
		"UserName": user.UserName,
		"NickName": user.NickName,
		"Email":    user.Email,
		"Status":   user.Status,
		"Avatar":   user.Avatar,
		"CreateAt": user.CreatedAt.Format(time.DateTime),
		"UpdateAt": user.UpdatedAt.Format(time.DateTime),
	}
}
