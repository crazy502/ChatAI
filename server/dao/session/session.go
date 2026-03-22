package session

import (
	"strings"
	"time"

	"server/common/mysql"
	"server/model"
)

func ListSessionsByUserName(userName, keyword string, includeArchived bool) ([]model.Session, error) {
	var sessions []model.Session

	query := mysql.DB.Model(&model.Session{}).Where("user_name = ?", userName)
	if !includeArchived {
		query = query.Where("archived = ?", false)
	}

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	err := query.
		Order("pinned desc").
		Order("last_message_at desc").
		Order("updated_at desc").
		Order("created_at desc").
		Find(&sessions).
		Error

	return sessions, err
}

func CreateSession(session *model.Session) (*model.Session, error) {
	err := mysql.DB.Create(session).Error
	return session, err
}

func GetSessionByID(sessionID string) (*model.Session, error) {
	var session model.Session
	err := mysql.DB.Where("id = ?", sessionID).First(&session).Error
	return &session, err
}

func GetSessionByIDAndUserName(sessionID, userName string) (*model.Session, error) {
	var session model.Session
	err := mysql.DB.Where("id = ? AND user_name = ?", sessionID, userName).First(&session).Error
	return &session, err
}

func UpdateSessionTitle(sessionID, userName, title string) error {
	return mysql.DB.Model(&model.Session{}).
		Where("id = ? AND user_name = ?", sessionID, userName).
		Updates(map[string]interface{}{
			"title":      title,
			"updated_at": time.Now(),
		}).
		Error
}

func UpdateSessionPin(sessionID, userName string, pinned bool) error {
	return mysql.DB.Model(&model.Session{}).
		Where("id = ? AND user_name = ?", sessionID, userName).
		Updates(map[string]interface{}{
			"pinned":     pinned,
			"updated_at": time.Now(),
		}).
		Error
}

func UpdateSessionArchive(sessionID, userName string, archived bool) error {
	return mysql.DB.Model(&model.Session{}).
		Where("id = ? AND user_name = ?", sessionID, userName).
		Updates(map[string]interface{}{
			"archived":   archived,
			"updated_at": time.Now(),
		}).
		Error
}

func TouchSession(sessionID string, lastMessageAt time.Time) error {
	return mysql.DB.Model(&model.Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_message_at": lastMessageAt,
			"updated_at":      lastMessageAt,
		}).
		Error
}
