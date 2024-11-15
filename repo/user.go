package repo

import (
	"log"

	"github.com/Peeranut-Kit/go_backend_test/utils"
	"gorm.io/gorm"
)

// Secondary port
type UserRepositoryInterface interface {
	CreateUser(user *utils.User) error
	GetUserFromEmail(user *utils.User) (*utils.User, error)
}

// Secondary adapter
type UserGormRepo struct {
	db *gorm.DB
}

// Initiate secondary adapter
func NewUserGormRepo(db *gorm.DB) UserRepositoryInterface {
	return &UserGormRepo{db: db}
}

func (r *UserGormRepo) CreateUser(user *utils.User) error {
	result := r.db.Create(user)

	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	return nil
}

func (r *UserGormRepo) GetUserFromEmail(user *utils.User) (*utils.User, error) {
	selectedUser := new(utils.User)
	result := r.db.Where("email = ?", user.Email).First(selectedUser)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return selectedUser, nil
}