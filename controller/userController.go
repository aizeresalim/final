package controller

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/aizeresalim/final/db"
	"github.com/aizeresalim/final/structures"
	"github.com/aizeresalim/final/tools"
	"gorm.io/gorm"
	"strconv"
)

// DeleteUser deletes a user account.
func DeleteUser(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	userID, err := tools.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Check if the user exists
	var user structures.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "User not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	if err := db.DB.Where("user_id = ?", userID).Delete(&structures.Blog{}).Error; err != nil {
		return c.JSON(fiber.Map{
			"message": "Failed to delete associated blogs",
		})
	}

	// Delete user from db
	if err := db.DB.Delete(&user).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User account deleted successfully",
	})
}

func UpdateUser(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	userID, err := tools.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Parse request body to extract updated user data
	var updatedUser structures.User
	if err := c.BodyParser(&updatedUser); err != nil {
		fmt.Println("Unable to parse user body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid user payload",
		})
	}

	// Check if the user exists
	var user structures.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "User not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Update user information
	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.Email = updatedUser.Email
	user.Phone = updatedUser.Phone

	// Save updated user to db
	if err := db.DB.Save(&user).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User information updated successfully",
		"user":    user,
	})
}

func GetUserInfo(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	userID, err := tools.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Query the db for user information
	var user structures.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "User not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Return user information
	return c.JSON(user)
}

// FollowUser allows a user to follow another user
func FollowUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	followerIDStr, err := tools.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Convert followerID to a uint
	followerID, err := strconv.ParseUint(followerIDStr, 10, 64)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Parse followed user ID from request parameters
	followedUserID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	// Check if the user is trying to follow themselves
	if followerID == followedUserID {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Cannot follow yourself",
		})
	}
	// Check if the followed user exists
	var followedUser structures.User
	if err := db.DB.First(&followedUser, followedUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Followed user not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Check if the follow relationship already exists
	var followRelationship structures.Follow
	if err := db.DB.Where("follower_id = ? AND followed_user_id = ?", followerID, followedUserID).First(&followRelationship).Error; err == nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Already following this user",
		})
	}

	// Create a new follow relationship
	follow := structures.Follow{
		FollowerID:     uint(followerID),
		FollowedUserID: uint(followedUserID),
	}
	if err := db.DB.Create(&follow).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to follow user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully followed user",
	})
}

// UnfollowUser allows a user to unfollow another user
func UnfollowUser(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	followerID, err := tools.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Parse followed user ID from request parameters
	followedUserID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	// Check if the follow relationship exists
	var followRelationship structures.Follow
	if err := db.DB.Where("follower_id = ? AND followed_user_id = ?", followerID, followedUserID).First(&followRelationship).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Not following this user",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Delete the follow relationship
	if err := db.DB.Delete(&followRelationship).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to unfollow user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully unfollowed user",
	})
}

// In controller/userController.go

// GetPostsFromFollowedUsers retrieves posts from users that the current user is following
func GetPostsFromFollowedUsers(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	userID, err := tools.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Get list of users the current user is following
	var followedUsers []structures.Follow
	if err := db.DB.Where("follower_id = ?", userID).Find(&followedUsers).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to retrieve followed users",
		})
	}

	// Extract followed user IDs
	var followedUserIDs []uint
	for _, follow := range followedUsers {
		followedUserIDs = append(followedUserIDs, follow.FollowedUserID)
	}

	// Retrieve posts from followed users
	var blogs []structures.Blog
	if err := db.DB.Where("user_id IN (?)", followedUserIDs).Find(&blogs).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to retrieve posts from followed users",
		})
	}

	return c.JSON(fiber.Map{
		"posts": blogs,
	})
}
