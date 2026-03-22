package webhook

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
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
	default:
		h.logger.Warnw("unknown or ignored webhook method", "method", req.FullMethod)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ignored"})
	}
}

func (h *handler) handleUserCreated(c *fiber.Ctx, req WebhookRequest) error {
	var resp AddHumanUserResponse
	if err := json.Unmarshal(req.Response, &resp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid response payload"})
	}

	zitadelUserID := strings.TrimSpace(resp.UserID)

	if zitadelUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	h.logger.Infow("processing user created webhook", "zitadel_user_id", zitadelUserID)

	account := models.Account{
		UserID:   &zitadelUserID,
		Type:     models.AccountTypeUser,
		Currency: "KZT",
	}

	_, err := h.accountDomain.CreateUsecase.Execute(c.UserContext(), account)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create account"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
