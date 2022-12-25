package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func GetLikedStatus(c *fiber.Ctx) error {
	var popularity models.Popularity
	postId := c.Params("postId")
	userId := c.Params("userId")
	var condition = map[string]interface{}{"post_id": postId, "user_id": userId}
	res := database.DB.Where(condition).Limit(1).Find(&popularity)
	if err := res.Error; err != nil {
		utils.ErrorResponse(c, utils.GetError)
	}
	if popularity.UserId == 0 {
		utils.ErrorResponse(c, utils.RecordNotFound)
	}

	return utils.GetRequestResponse(c, popularity)
}
