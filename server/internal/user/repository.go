package user

import (
	"errors"

	"server/infra/db"

	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	entity := new(User)
	err := db.Reader().Where("username = ?", username).First(entity).Error
	return entity, err
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	entity := new(User)
	err := db.Reader().Where("email = ?", email).First(entity).Error
	return entity, err
}

func (r *Repository) GetByEmailConsistent(email string) (*User, error) {
	entity := new(User)
	err := db.Writer().Where("email = ?", email).First(entity).Error
	return entity, err
}

func (r *Repository) Create(username, email, passwordHash string, isAdmin bool) (*User, error) {
	entity := &User{
		Email:    email,
		Name:     username,
		Username: username,
		Password: passwordHash,
		IsAdmin:  isAdmin,
	}
	return entity, db.Writer().Create(entity).Error
}

func (r *Repository) UpdatePassword(userID int64, passwordHash string) error {
	return db.Writer().Model(&User{}).
		Where("id = ?", userID).
		Update("password", passwordHash).
		Error
}

func (r *Repository) EnsureConfiguredAdmin(username, email, passwordHash string) error {
	return db.Writer().Transaction(func(tx *gorm.DB) error {
		var adminUser User
		err := tx.Where("username = ?", username).First(&adminUser).Error
		if err == gorm.ErrRecordNotFound {
			err = tx.Where("email = ?", email).First(&adminUser).Error
		}
		if err == gorm.ErrRecordNotFound {
			if passwordHash == "" {
				return errors.New("admin password is required when creating the initial administrator")
			}
			if err := tx.Model(&User{}).
				Where("username <> ?", username).
				Update("is_admin", false).Error; err != nil {
				return err
			}

			return tx.Create(&User{
				Email:    email,
				Name:     username,
				Username: username,
				Password: passwordHash,
				IsAdmin:  true,
			}).Error
		}
		if err != nil {
			return err
		}

		previousUsername := adminUser.Username

		if err := tx.Model(&User{}).
			Where("id <> ?", adminUser.ID).
			Update("is_admin", false).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"is_admin": true,
			"name":     username,
			"username": username,
		}
		if adminUser.Email == "" || adminUser.Email != email {
			updates["email"] = email
		}

		if err := tx.Model(&User{}).
			Where("id = ?", adminUser.ID).
			Updates(updates).
			Error; err != nil {
			return err
		}

		if previousUsername != "" && previousUsername != username {
			if err := tx.Table("sessions").
				Where("user_name = ?", previousUsername).
				Update("user_name", username).
				Error; err != nil {
				return err
			}

			if err := tx.Table("messages").
				Where("user_name = ?", previousUsername).
				Update("user_name", username).
				Error; err != nil {
				return err
			}
		}

		return nil
	})
}
