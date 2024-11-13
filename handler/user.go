package handler

import (
	"log"
	"os"
	"time"

	"github.com/Peeranut-Kit/go_backend_test/repo"
	"github.com/Peeranut-Kit/go_backend_test/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandlerInterface interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	GetCurrentUser(c *fiber.Ctx) error
}

// Primary adapter
type HttpUserHandler struct {
	UserRepo repo.UserRepositoryInterface
	validate *validator.Validate
}

// Initiate primary adapter
func NewHttpUserHandler(repo repo.UserRepositoryInterface, validate *validator.Validate) *HttpUserHandler {
	return &HttpUserHandler{UserRepo: repo, validate: validate}
}

func (u HttpUserHandler) Register(c *fiber.Ctx) error {
	user := new(utils.User)
	if err := c.BodyParser(user); err != nil {
		log.Println("Error decoding request body:", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// validate the user struct input
	if err := u.validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	// re-assign user password before saving in database
	user.Password = string(hashedPassword)

	err = u.UserRepo.CreateUser(user)
	if err != nil {
		log.Println("Error creating user:", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Create User Successful",
	})
}

func (u HttpUserHandler) Login(c *fiber.Ctx) error {
	user := new(utils.User)
	if err := c.BodyParser(user); err != nil {
		log.Println("Error decoding request body:", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// get user from email
	selectedUserByEmail, err := u.UserRepo.GetUserFromEmail(user)
	if err != nil {
		if err == utils.ErrNotFound {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(selectedUserByEmail.Password), []byte(user.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	// JWT part: Create the Claims
	claims := jwt.MapClaims{
		"user_id": selectedUserByEmail.ID,
		"name":    selectedUserByEmail.Name,
		"admin":   true,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response. (t is token)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Insert JWT token into Fiber Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"message": "Login success",
		"token":   t,
	})
}

func (u HttpUserHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	name := c.Locals("name").(string)

	return c.JSON(fiber.Map{
		"userID": userID,
		"name":   name,
	})
}
