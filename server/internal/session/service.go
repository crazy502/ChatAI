package session

import (
	"context"
	"strings"
	"unicode/utf8"

	"server/internal/ai"
	"server/pkg/apperror"
	"server/pkg/code"

	"gorm.io/gorm"
)

const (
	defaultSessionTitle  = "\u65b0\u4f1a\u8bdd"
	maxSessionTitleRunes = 100
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListByUserName(ctx context.Context, userName, keyword string, includeArchived bool) ([]SessionInfo, error) {
	sessions, err := s.repo.ListByUserName(userName, keyword, includeArchived)
	if err != nil {
		return nil, apperror.Wrap(code.CodeServerBusy, err, "list sessions failed").
			WithField("user_name", userName).
			WithField("include_archived", includeArchived)
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

func (s *Service) Rename(ctx context.Context, userName, sessionID, title string) error {
	if _, err := s.loadOwnedSession(ctx, userName, sessionID); err != nil {
		return err
	}

	if err := s.repo.UpdateTitle(sessionID, userName, NormalizeTitle(title)); err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "rename session failed").
			WithField("user_name", userName).
			WithField("session_id", sessionID)
	}

	return nil
}

func (s *Service) SetPinned(ctx context.Context, userName, sessionID string, pinned bool) error {
	if _, err := s.loadOwnedSession(ctx, userName, sessionID); err != nil {
		return err
	}

	if err := s.repo.UpdatePin(sessionID, userName, pinned); err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "update session pin failed").
			WithField("user_name", userName).
			WithField("session_id", sessionID).
			WithField("pinned", pinned)
	}

	return nil
}

func (s *Service) SetArchived(ctx context.Context, userName, sessionID string, archived bool) error {
	if _, err := s.loadOwnedSession(ctx, userName, sessionID); err != nil {
		return err
	}

	if err := s.repo.UpdateArchive(sessionID, userName, archived); err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "update session archive failed").
			WithField("user_name", userName).
			WithField("session_id", sessionID).
			WithField("archived", archived)
	}

	if archived {
		ai.GetGlobalManager().RemoveHelper(userName, sessionID)
	}

	return nil
}

func (s *Service) loadOwnedSession(ctx context.Context, userName, sessionID string) (*Session, error) {
	sessionInfo, err := s.repo.GetByIDAndUserName(sessionID, userName)
	if err == gorm.ErrRecordNotFound {
		return nil, apperror.New(code.CodeRecordNotFound, code.CodeRecordNotFound.Msg()).
			WithField("user_name", userName).
			WithField("session_id", sessionID)
	}
	if err != nil {
		return nil, apperror.Wrap(code.CodeServerBusy, err, "load session failed").
			WithField("user_name", userName).
			WithField("session_id", sessionID)
	}

	return sessionInfo, nil
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
