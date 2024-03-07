package controller

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/aizeresalim/final/db"
	"github.com/aizeresalim/final/structures"
	"github.com/aizeresalim/final/tools"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z0-9._%+\-]`)
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var userData structures.User

	// Parse form values
	email := c.FormValue("email")
	password := c.FormValue("password")
	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	phone := c.FormValue("phone")

	// Validate email format
	if len(email) == 0 || !validateEmail(strings.TrimSpace(email)) {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Email Address",
		})
	}

	// Validate password length
	if len(password) <= 6 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Password must be greater than 6 characters",
		})
	}

	// Check if email already exists
	db.DB.Where("email=?", email).First(&userData)
	if userData.Id != 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Create user
	user := structures.User{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Email:     email,
	}
	user.SetPassword(password)
	err := db.DB.Create(&user)
	if err != nil {
		log.Println(err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "creating user",
		})
	}

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Account created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	// Retrieve email and password from form data
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Validate email and password
	if email == "" || password == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Email and password are required",
		})
	}

	// Query the db to find the user by email
	var user structures.User
	db.DB.Where("email=?", email).First(&user)

	// Check if a user with the provided email exists
	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Email address doesn't exist, kindly create an account",
		})
	}

	// Compare the provided password with the hashed password stored in the db
	if err := user.ComparePassword(password); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	// Generate JWT token upon successful login
	token, err := tools.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	// Set the JWT token as an HTTP-only cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	// Return success message along with user data
	return c.JSON(fiber.Map{
		"message": "You have successfully logged in",
		"user":    user,
	})
}

type Claims struct {
	jwt.StandardClaims
}
