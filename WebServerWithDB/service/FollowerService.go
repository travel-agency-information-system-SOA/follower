package service

import (
	"database-example/model"
	"database-example/repo"
	"fmt"
)

type FollowerService struct {
	FollowerRepo *repo.FollowerRepository
}

func (service *FollowerService) FindFollower(id string) (*model.Followers, error) {
	tour, err := service.FollowerRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
	}
	return &tour, nil
}

func (service *FollowerService) Create(follower *model.Followers) error {
	err := service.FollowerRepo.CreateFollower(follower)
	if err != nil {
		return err
	}
	return nil
}
