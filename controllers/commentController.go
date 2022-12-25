package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Comments(c *fiber.Ctx) error {
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

	user_id, err := utils.ParseUint(data["user_id"])
	if err != nil {
		return err
	}
	postId, err := utils.ParseUint(c.Params("postId"))
	if err != nil {
		return utils.ErrorResponse(c, utils.GetError)
	}

	comment := models.Comment{
		UserId:    user_id,
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
	commentId := c.Params("commentId")
	if err := database.DB.Delete(&models.Comment{}, commentId).Error; err != nil {
		return utils.ErrorResponse(c, utils.DeleteError)
	}

	return utils.ResponseBody(c, utils.DeleteSuccess)
}
