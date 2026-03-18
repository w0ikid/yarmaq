package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
	"github.com/w0ikid/yarmaq/pkg/zitadel"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"go.uber.org/zap"
)

type Service interface {
	// our repo methods
	Create(ctx context.Context, user models.User) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByZitadelID(ctx context.Context, zitadelID string) (*models.User, error)
	Update(ctx context.Context, user models.User) (*models.User, error)
	// zitadel cleint methods
	GetFromZitadel(ctx context.Context, zitadelID string) (*management.GetUserByIDResponse, error)
	AssignRole(ctx context.Context, zitadelID string, roles []string) error
}

type implementation struct {
	repo    UsersRepo
	zitadel *zitadel.Client
	logger  *zap.SugaredLogger
}

func NewService(repo UsersRepo, zitadelCleint *zitadel.Client, logger *zap.SugaredLogger) Service {
	return &implementation{
		repo:    repo,
		zitadel: zitadelCleint,
		logger:  logger.Named("users_service"),
	}
}

// CREATE
func (s *implementation) Create(ctx context.Context, user models.User) (*models.User, error) {
	s.logger.Infow("creating user", "email", user.Email)
	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Errorw("failed to create user", "email", user.Email, "error", err)
		return nil, err
	}
	s.logger.Infow("user created successfully", "id", createdUser.ID)
	return createdUser, nil
}

// GET BY ID
func (s *implementation) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	s.logger.Infow("getting user by id", "id", id)
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("failed to get user by id", "id", id, "error", err)
		return nil, err
	}
	if user == nil {
		s.logger.Infow("user not found by id", "id", id)
		return nil, nil
	}
	s.logger.Infow("user found by id", "id", id)
	return user, nil
}

// GET BY EMAIL
func (s *implementation) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	s.logger.Infow("getting user by email", "email", email)
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Errorw("failed to get user by email", "email", email, "error", err)
		return nil, err
	}
	if user == nil {
		s.logger.Infow("user not found by email", "email", email)
		return nil, nil
	}
	s.logger.Infow("user found by email", "email", email)
	return user, nil
}

// GET BY ZITADEL ID
func (s *implementation) GetByZitadelID(ctx context.Context, zitadelID string) (*models.User, error) {
	s.logger.Infow("getting user by zitadel id", "zitadel_id", zitadelID)
	user, err := s.repo.GetByZitadelID(ctx, zitadelID)
	if err != nil {
		s.logger.Errorw("failed to get user by zitadel id", "zitadel_id", zitadelID, "error", err)
		return nil, err
	}
	if user == nil {
		s.logger.Infow("user not found by zitadel id", "zitadel_id", zitadelID)
		return nil, nil
	}
	s.logger.Infow("user found by zitadel id", "zitadel_id", zitadelID)
	return user, nil
}

// UPDATE
func (s *implementation) Update(ctx context.Context, user models.User) (*models.User, error) {
	s.logger.Infow("updating user", "id", user.ID)
	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		s.logger.Errorw("failed to update user", "id", user.ID, "error", err)
		return nil, err
	}
	s.logger.Infow("user updated successfully", "id", updatedUser.ID)
	return updatedUser, nil
}

func (s *implementation) UpdateRoles(ctx context.Context, zitadelID string, roles []string) (*models.User, error) {
	s.logger.Infow("updating roles", "zitadel_id", zitadelID, "roles", roles)
	user, err := s.repo.UpdateRoles(ctx, zitadelID, roles)
	if err != nil {
		s.logger.Errorw("failed to update roles", "zitadel_id", zitadelID, "error", err)
		return nil, err
	}
	s.logger.Infow("roles updated", "zitadel_id", zitadelID)
	return user, nil
}

// ZITADEL CLIENT METHODS
func (s *implementation) AssignRole(ctx context.Context, zitadelID string, roles []string) error {
	s.logger.Infow("assigning roles", "zitadel_id", zitadelID, "roles", roles)
	_, err := s.zitadel.Mgmt.AddUserGrant(ctx, &management.AddUserGrantRequest{
		UserId:    zitadelID,
		ProjectId: "364528156301852676",
		RoleKeys:  roles,
	})
	if err != nil {
		s.logger.Errorw("failed to assign roles", "zitadel_id", zitadelID, "error", err)
		return err
	}
	s.logger.Infow("roles assigned", "zitadel_id", zitadelID)
	return nil
}

func (s *implementation) GetFromZitadel(ctx context.Context, zitadelID string) (*management.GetUserByIDResponse, error) {
	s.logger.Infow("getting user from zitadel", "zitadel_id", zitadelID)
	user, err := s.zitadel.Mgmt.GetUserByID(ctx, &management.GetUserByIDRequest{
		Id: zitadelID,
	})
	if err != nil {
		s.logger.Errorw("failed to get user from zitadel", "zitadel_id", zitadelID, "error", err)
		return nil, err
	}
	s.logger.Infow("user found in zitadel", "zitadel_id", zitadelID)
	return user, nil
}
