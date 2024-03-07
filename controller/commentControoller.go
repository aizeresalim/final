package controller

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/aizeresalim/final/db"
	"github.com/aizeresalim/final/structures"
	"github.com/aizeresalim/final/tools"
	"gorm.io/gorm"
)

func CreateComment(c *fiber.Ctx) error {
	// Parse request body to extract comment data
	var commentData structures.Comment
	if err := c.BodyParser(&commentData); err != nil {
		fmt.Println("Unable to parse comment body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid comment payload",
		})
	}

	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	userID, err := tools.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Parse blog post ID from URL parameter
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}

	// Check if the blog post exists
	var blogPost structures.Blog
	if err := db.DB.Where("id = ?", postID).First(&blogPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Blog post not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Create new comment object
	comment := structures.Comment{
		UserID:   userID,
		PostID:   uint(postID),
		Content:  commentData.Content,
		DateTime: time.Now(),
	}

	// Save comment to db
	if err := db.DB.Create(&comment).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to create comment",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// UpdateComment updates an existing comment.
func UpdateComment(c *fiber.Ctx) error {
	// Parse comment ID from URL parameter
	commentID, err := strconv.Atoi(c.Params("commentID"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid comment ID",
		})
	}

	// Parse request body to extract updated comment data
	var updatedComment structures.Comment
	if err := c.BodyParser(&updatedComment); err != nil {
		fmt.Println("Unable to parse comment body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid comment payload",
		})
	}

	// Check if the comment exists
	var comment structures.Comment
	if err := db.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Comment not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Update comment content
	comment.Content = updatedComment.Content

	// Save updated comment to db
	if err := db.DB.Save(&comment).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to update comment",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comment updated successfully",
		"comment": comment,
	})
}

// DeleteComment deletes an existing comment.
func DeleteComment(c *fiber.Ctx) error {
	// Parse comment ID from URL parameter
	commentID, err := strconv.Atoi(c.Params("commentID"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid comment ID",
		})
	}

	// Check if the comment exists
	var comment structures.Comment
	if err := db.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Comment not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Delete comment from db
	if err := db.DB.Delete(&comment).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to delete comment",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comment deleted successfully",
	})
}

func ReadComments(c *fiber.Ctx) error {
	// Parse blog post ID from URL parameter
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}

	// Check if the blog post exists
	var blogPost structures.Blog
	if err := db.DB.Where("id = ?", postID).First(&blogPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Blog post not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Retrieve comments associated with the blog post
	var comments []structures.Comment
	if err := db.DB.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to retrieve comments",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Comments retrieved successfully",
		"comments": comments,
	})
}
