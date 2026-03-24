package user

import (
	"server/infra/db"

	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	entity := new(User)
	err := db.DB.Where("username = ?", username).First(entity).Error
	return entity, err
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	entity := new(User)
	err := db.DB.Where("email = ?", email).First(entity).Error
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
	return entity, db.DB.Create(entity).Error
}

func (r *Repository) UpdatePassword(userID int64, passwordHash string) error {
	return db.DB.Model(&User{}).
		Where("id = ?", userID).
		Update("password", passwordHash).
		Error
}

func (r *Repository) EnsureConfiguredAdmin(username, email, passwordHash string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&User{}).
			Where("username <> ?", username).
			Update("is_admin", false).Error; err != nil {
			return err
		}

		var adminUser User
		err := tx.Where("username = ?", username).First(&adminUser).Error
		if err == gorm.ErrRecordNotFound {
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

		updates := map[string]interface{}{
			"is_admin": true,
			"name":     username,
			"password": passwordHash,
		}
		if adminUser.Email == "" || adminUser.Email != email {
			updates["email"] = email
		}

		return tx.Model(&User{}).
			Where("id = ?", adminUser.ID).
			Updates(updates).
			Error
	})
}
