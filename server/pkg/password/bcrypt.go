package password

import "golang.org/x/crypto/bcrypt"

// HashPassword 对密码进行哈希处理
func HashPassword(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword 验证密码是否匹配
func CheckPassword(hash, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}

// IsBcryptHash 验证哈希是否为bcrypt哈希
func IsBcryptHash(hash string) bool {
	_, err := bcrypt.Cost([]byte(hash))
	return err == nil
}
