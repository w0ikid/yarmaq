package webhook

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type WebhookRequest struct {
	FullMethod string              `json:"fullMethod"`
	InstanceID string              `json:"instanceID"`
	OrgID      string              `json:"orgID"`
	ProjectID  string              `json:"projectID"`
	UserID     string              `json:"userID"`
	Request    json.RawMessage     `json:"request"`
	Response   json.RawMessage     `json:"response"`
	Headers    map[string][]string `json:"headers"`
}

type AddHumanUserRequest struct {
	Username string `json:"username"`
	Profile  struct {
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	} `json:"profile"`
	Email struct {
		Email      string `json:"email"`
		IsVerified bool   `json:"isVerified"`
	} `json:"email"`
}

type AddHumanUserResponse struct {
	UserID string `json:"userId"`
}

type UpdateUserGrantRequest struct {
	UserID   string   `json:"userId"`
	GrantID  string   `json:"grantId"`
	RoleKeys []string `json:"roleKeys"`
}

func (h *handler) HandleZitadelSync(c *fiber.Ctx) error {
	var req WebhookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.FullMethod == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fullMethod"})
	}

	switch req.FullMethod {
	case "/zitadel.user.v2.UserService/AddHumanUser":
		return h.handleUserCreated(c, req)
	case "/zitadel.management.v1.ManagementService/UpdateUserGrant":
		return h.handleRoleUpdated(c, req)
	default:
		h.logger.Warnw("unknown webhook method", "method", req.FullMethod)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ignored"})
	}
}

func (h *handler) handleUserCreated(c *fiber.Ctx, req WebhookRequest) error {
	var r AddHumanUserRequest
	if err := json.Unmarshal(req.Request, &r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	var resp AddHumanUserResponse
	if err := json.Unmarshal(req.Response, &resp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid response payload"})
	}

	email := strings.TrimSpace(r.Email.Email)
	username := strings.TrimSpace(r.Username)
	if username == "" {
		username = email
	}
	userID := strings.TrimSpace(resp.UserID)

	if userID == "" || email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	h.logger.Infow("user created webhook", "user_id", userID, "email", email)

	user := models.User{
		ZitadelUserID: userID,
		Email:         email,
		Username:      username,
		Roles:         []string{"user"},
		IsActive:      true,
	}

	createdUser, err := h.createUser.Execute(c.UserContext(), user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			h.logger.Infow("user already exists", "zitadel_user_id", userID, "email", email)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "message": "user already exists"})
		}
		h.logger.Errorw("failed to create user", "zitadel_user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	h.logger.Infow("user created", "id", createdUser.ID, "zitadel_user_id", createdUser.ZitadelUserID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "ok", "user": createdUser})
}

func (h *handler) handleRoleUpdated(c *fiber.Ctx, req WebhookRequest) error {
	var r UpdateUserGrantRequest
	if err := json.Unmarshal(req.Request, &r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	userID := strings.TrimSpace(r.UserID)
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing userId"})
	}

	// r.RoleKeys может быть пустым — это значит все роли убраны
	h.logger.Infow("role update webhook", "user_id", userID, "roles", r.RoleKeys)

	if _, err := h.updateUser.Execute(c.UserContext(), models.User{
		ZitadelUserID: userID,
		Roles:         r.RoleKeys, // пустой slice — очищает роли
	}); err != nil {
		h.logger.Errorw("failed to update roles", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update roles"})
	}

	h.logger.Infow("roles updated", "user_id", userID, "roles", r.RoleKeys)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
