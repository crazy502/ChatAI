package session

import (
	"strings"
	"time"

	"server/infra/db"
)

type Repository struct{}

const listQueryIndexName = "idx_sessions_user_archived_pinned_activity"

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) ListByUserName(userName, keyword string, includeArchived bool) ([]Session, error) {
	var sessions []Session

	query := db.Reader().Model(&Session{}).Where("user_name = ?", userName)
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

func (r *Repository) Create(entity *Session) (*Session, error) {
	return entity, db.Writer().Create(entity).Error
}

func (r *Repository) GetByID(sessionID string) (*Session, error) {
	entity := new(Session)
	err := db.Writer().Where("id = ?", sessionID).First(entity).Error
	return entity, err
}

func (r *Repository) GetByIDAndUserName(sessionID, userName string) (*Session, error) {
	entity := new(Session)
	err := db.Writer().Where("id = ? AND user_name = ?", sessionID, userName).First(entity).Error
	return entity, err
}

func (r *Repository) UpdateTitle(sessionID, userName, title string) error {
	return db.Writer().Model(&Session{}).
		Where("id = ? AND user_name = ?", sessionID, userName).
		Updates(map[string]interface{}{
			"title":      title,
			"updated_at": time.Now(),
		}).
		Error
}

func (r *Repository) UpdatePin(sessionID, userName string, pinned bool) error {
	return db.Writer().Model(&Session{}).
		Where("id = ? AND user_name = ?", sessionID, userName).
		Updates(map[string]interface{}{
			"pinned":     pinned,
			"updated_at": time.Now(),
		}).
		Error
}

func (r *Repository) UpdateArchive(sessionID, userName string, archived bool) error {
	return db.Writer().Model(&Session{}).
		Where("id = ? AND user_name = ?", sessionID, userName).
		Updates(map[string]interface{}{
			"archived":   archived,
			"updated_at": time.Now(),
		}).
		Error
}

func (r *Repository) TouchSession(sessionID string, lastMessageAt time.Time) error {
	return db.Writer().Model(&Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_message_at": lastMessageAt,
			"updated_at":      lastMessageAt,
		}).
		Error
}

func (r *Repository) EnsureListIndexes() error {
	writer := db.Writer()
	if writer.Migrator().HasIndex(&Session{}, listQueryIndexName) {
		return nil
	}

	return writer.Exec(
		"CREATE INDEX " + listQueryIndexName + " ON sessions (user_name, archived, pinned, last_message_at, updated_at, created_at)",
	).Error
}
