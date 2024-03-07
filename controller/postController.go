package controller

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/aizeresalim/final/db"
	"github.com/aizeresalim/final/structures"
	"github.com/aizeresalim/final/tools"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	// Retrieve the user ID from the session
	userID := c.Locals("userID").(string) // Assuming userID is stored as a string in the session

	// Parse the request body into a structures.Blog struct
	var blogpost structures.Blog
	if err := c.BodyParser(&blogpost); err != nil {
		fmt.Println("Unable to parse body")
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid payload",
		})
	}

	// Set the UserID field of the blogpost with the retrieved user ID
	blogpost.UserID = userID

	// Create the blog post in the db
	if err := db.DB.Create(&blogpost).Error; err != nil {
		fmt.Println("Error creating post:", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error creating post",
		})
	}

	// Return a success response if the blog post was created successfully
	return c.JSON(fiber.Map{
		"message": "Congratulations! Your post is live",
	})
}

func AllPost(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 5
	offset := (page - 1) * limit
	var total int64
	var getblog []structures.Blog
	db.DB.Preload("User").Offset(offset).Limit(limit).Find(&getblog)
	db.DB.Model(&structures.Blog{}).Count(&total)
	return c.JSON(fiber.Map{
		"data": getblog,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(int(total) / limit)),
		},
	})

}

func DetailPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var blogpost structures.Blog
	db.DB.Where("id=?", id).Preload("User").First(&blogpost)
	return c.JSON(fiber.Map{
		"data": blogpost,
	})

}

func UpdatePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	blog := structures.Blog{
		Id: uint(id),
	}

	if err := c.BodyParser(&blog); err != nil {
		fmt.Println("Unable to parse body")
	}
	db.DB.Model(&blog).Updates(blog)
	return c.JSON(fiber.Map{
		"message": "post updated successfully",
	})

}

func UniquePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, _ := tools.Parsejwt(cookie)
	var blog []structures.Blog
	db.DB.Model(&blog).Where("user_id=?", id).Preload("User").Find(&blog)

	return c.JSON(blog)

}
func DeletePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	blog := structures.Blog{
		Id: uint(id),
	}
	deleteQuery := db.DB.Delete(&blog)
	if errors.Is(deleteQuery.Error, gorm.ErrRecordNotFound) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Opps!, record Not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "post deleted successfully",
	})

}
