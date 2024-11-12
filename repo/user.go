package repo

import (
	"log"

	"github.com/Peeranut-Kit/go_backend_test/utils"
)

func (postgres *PostgresDB) CreateUser(user *utils.User) error {
	result := postgres.db.Create(user)

	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	return nil
}

func (postgres *PostgresDB) GetUserFromEmail(user *utils.User) (*utils.User, error) {
	selectedUser := new(utils.User)
	result := postgres.db.Where("email = ?", user.Email).First(selectedUser)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return selectedUser, nil
}