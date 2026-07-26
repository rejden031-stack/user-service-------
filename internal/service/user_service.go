package service

import (
	"context"
	"errors"
	"log"

	"user-service/internal/entity"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Create(ctx context.Context, name string) (entity.User, error)
}

type EventPublisher interface {
	PublishUserCreated(ctx context.Context, id int, name string) error
}

type UserService struct {
	repo      UserRepository
	publisher EventPublisher
}

func NewUserService(repo UserRepository, publisher EventPublisher) *UserService {
	return &UserService{repo: repo, publisher: publisher}
}

func (s *UserService) GetUser(ctx context.Context, id int) (entity.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, name string) (entity.User, error) {
	if name == "" {
		return entity.User{}, errors.New("name is required")
	}

	user, err := s.repo.Create(ctx, name)
	if err != nil {
		return entity.User{}, err
	}


	if err := s.publisher.PublishUserCreated(ctx, user.ID, user.Name); err != nil {
		log.Printf("failed to publish user created event: %v", err)
	}

	return user, nil
}
