package session

import (
	"log"
	"strings"
	"unicode/utf8"

	"server/internal/ai"
	"server/pkg/code"

	"gorm.io/gorm"
)

const (
	defaultSessionTitle  = "新会话"
	maxSessionTitleRunes = 100
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListByUserName(userName, keyword string, includeArchived bool) ([]SessionInfo, error) {
	sessions, err := s.repo.ListByUserName(userName, keyword, includeArchived)
	if err != nil {
		return nil, err
	}

	result := make([]SessionInfo, 0, len(sessions))
	for _, item := range sessions {
		result = append(result, SessionInfo{
			ID:              item.ID,
			SessionID:       item.ID,
			LegacySessionID: item.ID,
			Title:           item.Title,
			LegacyTitle:     item.Title,
			Pinned:          item.Pinned,
			Archived:        item.Archived,
			LastMessageAt:   item.LastMessageAt,
			UpdatedAt:       item.UpdatedAt,
		})
	}

	return result, nil
}

func (s *Service) Rename(userName, sessionID, title string) code.Code {
	if _, resultCode := s.loadOwnedSession(userName, sessionID); resultCode != code.CodeSuccess {
		return resultCode
	}

	if err := s.repo.UpdateTitle(sessionID, userName, NormalizeTitle(title)); err != nil {
		log.Println("rename session error:", err)
		return code.CodeServerBusy
	}

	return code.CodeSuccess
}

func (s *Service) SetPinned(userName, sessionID string, pinned bool) code.Code {
	if _, resultCode := s.loadOwnedSession(userName, sessionID); resultCode != code.CodeSuccess {
		return resultCode
	}

	if err := s.repo.UpdatePin(sessionID, userName, pinned); err != nil {
		log.Println("set session pinned error:", err)
		return code.CodeServerBusy
	}

	return code.CodeSuccess
}

func (s *Service) SetArchived(userName, sessionID string, archived bool) code.Code {
	if _, resultCode := s.loadOwnedSession(userName, sessionID); resultCode != code.CodeSuccess {
		return resultCode
	}

	if err := s.repo.UpdateArchive(sessionID, userName, archived); err != nil {
		log.Println("set session archived error:", err)
		return code.CodeServerBusy
	}

	if archived {
		ai.GetGlobalManager().RemoveHelper(userName, sessionID)
	}

	return code.CodeSuccess
}

func (s *Service) loadOwnedSession(userName, sessionID string) (*Session, code.Code) {
	sessionInfo, err := s.repo.GetByIDAndUserName(sessionID, userName)
	if err == gorm.ErrRecordNotFound {
		return nil, code.CodeRecordNotFound
	}
	if err != nil {
		log.Println("load owned session error:", err)
		return nil, code.CodeServerBusy
	}

	return sessionInfo, code.CodeSuccess
}

func NormalizeTitle(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return defaultSessionTitle
	}

	if utf8.RuneCountInString(title) <= maxSessionTitleRunes {
		return title
	}

	runes := []rune(title)
	return strings.TrimSpace(string(runes[:maxSessionTitleRunes]))
}
