package user

import (
	"server/common/mysql"
	"server/model"

	"gorm.io/gorm"
)

const (
	CodeMsg     = "GopherAI楠岃瘉鐮佸涓?楠岃瘉鐮佷粎闄愪簬2鍒嗛挓鏈夋晥): "
	UserNameMsg = "GopherAI鐨勮处鍙峰涓嬶紝璇蜂繚鐣欏ソ锛屽悗缁彲浠ョ敤璐﹀彿杩涜鐧诲綍 "
)

// 杩欒竟鍙兘閫氳繃璐﹀彿杩涜鐧诲綍
func IsExistUser(username string) (bool, *model.User) {
	user, err := mysql.GetUserByUsername(username)
	if err == gorm.ErrRecordNotFound || user == nil {
		return false, nil
	}

	return true, user
}

// 妫€鏌ラ偖绠辨槸鍚﹀凡瀛樺湪
func IsExistEmail(email string) (bool, *model.User) {
	user, err := mysql.GetUserByEmail(email)
	if err == gorm.ErrRecordNotFound || user == nil {
		return false, nil
	}

	return true, user
}

func Register(username, email, passwordHash string, isAdmin bool) (*model.User, bool) {
	if user, err := mysql.InsertUser(&model.User{
		Email:    email,
		Name:     username,
		Username: username,
		Password: passwordHash,
		IsAdmin:  isAdmin,
	}); err != nil {
		return nil, false
	} else {
		return user, true
	}
}

func HasAdminUsers() (bool, error) {
	var count int64
	if err := mysql.DB.Model(&model.User{}).
		Where("is_admin = ?", true).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func UpdatePassword(userID int64, passwordHash string) error {
	return mysql.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Update("password", passwordHash).
		Error
}
