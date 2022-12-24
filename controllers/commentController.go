package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Comments(c *fiber.Ctx) error {
	// Check that user is logged in
	if _, err := utils.GetCurrentUser(c, SecretKey); err != nil {
		return utils.ErrorResponse(c, utils.UserNotFound)
	}

	comments := []models.Comment{}
	postId := c.Params("postId")
	condition := map[string]interface{}{"post_id": postId}

	if err := database.DB.Joins("User").Where(condition).Find(&comments).Error; err != nil {
		return utils.ErrorResponse(c, utils.GetError)
	}

	return utils.GetRequestResponse(c, comments)
}

func AddComment(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Check that user is logged in
	user, err := utils.GetCurrentUser(c, SecretKey)
	if err != nil {
		return utils.ErrorResponse(c, utils.UserNotFound)
	}
	postId, err := utils.ParseUint(c.Params("postId"))
	if err != nil {
		return utils.ErrorResponse(c, utils.GetError)
	}

	comment := models.Comment{
		UserId:    user.UserId,
		PostId:    postId,
		Content:   data["content"],
		CreatedDt: time.Now().Unix(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		return utils.ErrorResponse(c, utils.CreateError)
	}

	return utils.CreateRequestResponse(c, comment)
}

func DeleteComment(c *fiber.Ctx) error {
	// Check that user is logged in
	user, err := utils.GetCurrentUser(c, SecretKey)
	if err != nil {
		return utils.ErrorResponse(c, utils.UserNotFound)
	}

	// Check that user is an administrator
	if user.AccessType == 1 {
		commentId := c.Params("commentId")
		if err := database.DB.Delete(&models.Comment{}, commentId).Error; err != nil {
			return utils.ErrorResponse(c, utils.DeleteError)
		}

		return utils.ResponseBody(c, utils.DeleteSuccess)
	}

	return utils.ErrorResponse(c, utils.ForbiddenAction)
}
