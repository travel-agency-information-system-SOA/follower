package repo

import (
	"database-example/model"

	"gorm.io/gorm"
)

type FollowerRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *FollowerRepository) FindById(id string) (model.Followers, error) {
	follower := model.Followers{}
	dbResult := repo.DatabaseConnection.First(&follower, "id = ?", id)
	if dbResult != nil {
		return follower, dbResult.Error
	}
	return follower, nil
}

func (repo *FollowerRepository) CreateFollower(follower *model.Followers) error {
	dbResult := repo.DatabaseConnection.Create(follower)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
