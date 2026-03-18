package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/users"
	"github.com/w0ikid/yarmaq/pkg/ctxkeys"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	CreateUser users.CreateUsecase
	GetUser    users.GetUserUsecase
	Logger     *zap.SugaredLogger
}

type Handler interface {
	CreateUser(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	ProtectedUserEndpoint(c *fiber.Ctx) error
}

type handler struct {
	createUser users.CreateUsecase
	getUser    users.GetUserUsecase
	logger     *zap.SugaredLogger
}

func NewHandler(deps HandlerDeps) Handler {
	log := deps.Logger.Named("users_handler")
	return &handler{
		createUser: deps.CreateUser,
		getUser:    deps.GetUser,
		logger:     log,
	}
}

func (h *handler) CreateUser(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"message": "CreateUser endpoint"})
}

func (h *handler) GetUser(c *fiber.Ctx) error {
	userID := c.Locals("userID")

	user, err := h.getUser.Execute(c.Context(), userID.(string))
	if err != nil {
		h.logger.Errorw("error executing GetUser usecase", "userID", userID, "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": "An error occurred while fetching the user",
		})
	}
	if user == nil {
		h.logger.Infow("user not found", "userID", userID)
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "User not found",
		})
	}

	h.logger.Infow("successfully fetched user", "userID", userID)
	return c.Status(200).JSON(user)
}

func (h *handler) ProtectedUserEndpoint(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	roles, _ := c.Locals("roles").([]string)

	ctx := ctxkeys.WithUserContext(c.UserContext(), userID, roles)

	ctxkeys.GetUserID(ctx)
	ctxkeys.GetRoles(ctx)
	
	h.logger.Infow("accessed protected endpoint", "userID", userID)
	return c.Status(200).JSON(fiber.Map{"message": "This is a protected endpoint", "userID": userID})
}
