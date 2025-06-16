// pkg/internal/service/user/user.go
package user

import ()

var (
// ErrUserNotFound 用户不存在
// ErrUserNotFound = errors.New("user not found")
// ErrInvalidPassword 密码不正确
// ErrInvalidPassword = errors.New("invalid password")
)

// Authenticate 验证用户名/密码，成功返回 userID，失败返回错误
func Authenticate(username, password string) (uint, error) {
	// 1. 从仓库层（Repository）取出用户

	// 3. 验证通过，返回用户 ID
	return 0, nil
}
