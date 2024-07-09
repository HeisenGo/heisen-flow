// /*
// This middleware will check the user's permissions for each request.
// */

package middlewares

// import (
// 	userboard "server/internal/user_board"
// 	"server/pkg/rbac"

// 	"github.com/gofiber/fiber/v2"
// )

// func RBACMiddleware(userBoardOps userboard.Ops) func(requiredPermission rbac.Permission) fiber.Handler {
// 	return func(requiredPermission rbac.Permission) fiber.Handler {
// 		return func(c *fiber.Ctx) error {
// 			userID, ok := valuecontext.GetUserID(c.Context())
// 			if !ok {
// 				return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
// 			}

// 			boardID := c.Query("boardID")
// 			if boardID == "" {
// 				return c.Status(fiber.StatusBadRequest).SendString("Board ID is required")
// 			}

// 			role, err := userBoardOps.GetUserBoardRole(c.Context(), userID, boardID)
// 			if err != nil {
// 				return c.Status(fiber.StatusInternalServerError).SendString("Error getting user role")
// 			}

// 			if !rbac.HasPermission(role, requiredPermission) {
// 				return c.Status(fiber.StatusForbidden).SendString("Permission denied")
// 			}

// 			c.Locals(valuecontext.UserRoleKey, role)

// 			return c.Next()
// 		}
// 	}
// }
