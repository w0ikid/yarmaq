package webhook

import (
	"encoding/json"
	"fmt"
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
	zitadelUserID := strings.TrimSpace(resp.UserID)

	if zitadelUserID == "" || email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	h.logger.Infow("processing user created webhook", "zitadel_user_id", zitadelUserID, "email", email)

	account := models.Account{
		UserID:   resp.UserID,
		Number:   fmt.Sprintf("ACC-%s", strings.ToUpper(resp.UserID[:8])),
		Currency: "KZT",
		Status:   "ACTIVE",
	}

	_, err := h.accountDomain.CreateUsecase.Execute(c.UserContext(), account)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create account"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
