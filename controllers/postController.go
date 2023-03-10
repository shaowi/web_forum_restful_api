package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Posts(c *fiber.Ctx) error {
	posts := []models.Post{}
	if err := database.DB.Joins("User").Find(&posts).Error; err != nil {
		return utils.ErrorResponse(c, utils.GetError)
	}

	for i, post := range posts {
		posts[i] = utils.GetPostStats(post)
	}
	return utils.GetRequestResponse(c, posts)
}

func ViewSeenPosts(c *fiber.Ctx) error {
	userId, err := utils.ParseUint(c.Params("userId"))
	if err != nil {
		return err
	}

	condition := map[string]interface{}{"popularities.user_id": userId}
	joinPopularityCondition := "JOIN popularities ON popularities.post_id = posts.post_id"
	return utils.GetHistoryPosts(c, condition, joinPopularityCondition)
}

func ViewLikedPosts(c *fiber.Ctx) error {
	userId, err := utils.ParseUint(c.Params("userId"))
	if err != nil {
		return err
	}

	condition := map[string]interface{}{"popularities.user_id": userId, "popularities.likes": true}
	joinPopularityCondition := "JOIN popularities ON popularities.post_id = posts.post_id"
	return utils.GetHistoryPosts(c, condition, joinPopularityCondition)
}

func ViewCommentedPosts(c *fiber.Ctx) error {
	userId, err := utils.ParseUint(c.Params("userId"))
	if err != nil {
		return err
	}

	condition := map[string]interface{}{"comments.user_id": userId}
	joinCommentCondition := "JOIN comments ON comments.post_id = posts.post_id"
	return utils.GetHistoryPosts(c, condition, joinCommentCondition)
}

func AddPost(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	user_id, err := utils.ParseUint(data["user_id"])
	if err != nil {
		return err
	}

	post := models.Post{
		UserId:     user_id,
		Title:      data["title"],
		Body:       data["body"],
		Categories: data["categories"],
		CreatedDt:  time.Now().Unix(),
	}

	if err := database.DB.Create(&post).Error; err != nil {
		return utils.ErrorResponse(c, utils.CreateError)
	}

	if err := database.DB.Joins("User").Find(&post).Error; err != nil {
		return utils.ErrorResponse(c, utils.GetError)
	}

	return utils.CreateRequestResponse(c, post)
}

func DeletePost(c *fiber.Ctx) error {
	postId := c.Params("postId")
	if err := database.DB.Delete(&models.Post{}, postId).Error; err != nil {
		return utils.ErrorResponse(c, utils.DeleteError)
	}
	return utils.ResponseBody(c, utils.DeleteSuccess)

}

func LikePost(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Like: 1, Unlike: 0
	like, err := strconv.ParseBool(data["like"])
	if err != nil {
		return err
	}

	// Incrementing the likes of a post
	var post models.Post

	// Search for post
	postId, err := utils.ParseUint(c.Params("postId"))
	if err != nil {
		return err
	}
	var condition = map[string]interface{}{"post_id": postId}
	if err := database.DB.Where(condition).First(&post).Error; err != nil {
		return utils.ErrorResponse(c, utils.GetError)
	}
	var curLikes uint = post.Likes
	newLikeCnt := curLikes + 1

	if !like {
		if curLikes == 0 {
			newLikeCnt = curLikes
		} else {
			newLikeCnt = curLikes - 1
		}
	}

	// Update
	if err := database.DB.Model(&post).Update("likes", newLikeCnt).Error; err != nil {
		return utils.ResponseBody(c, utils.UpdateError)
	}

	// Update the like status of this post for current user based on type
	var popularity models.Popularity

	user_id, err := utils.ParseUint(data["user_id"])
	if err != nil {
		return err
	}

	// Search
	condition = map[string]interface{}{"post_id": postId, "user_id": user_id}
	res := database.DB.Where(condition).Limit(1).Find(&popularity)
	if res.RowsAffected == 0 {
		// Create new record as user info has not been recorded yet for this post
		popularity = models.Popularity{
			UserId: user_id,
			PostId: postId,
			Likes:  like,
		}
		if err := database.DB.Create(&popularity).Error; err != nil {
			return utils.ResponseBody(c, utils.CreateError)
		}
	} else {
		// Update like status
		if err := database.DB.Model(&popularity).Update("likes", like).Error; err != nil {
			return utils.ResponseBody(c, utils.UpdateError)
		}
	}

	return utils.ResponseBody(c, utils.Success)
}

func ViewPost(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	user_id, err := utils.ParseUint(data["user_id"])
	if err != nil {
		return err
	}

	// Incrementing the views of a post
	var post models.Post

	// Search for post
	postId, err := utils.ParseUint(c.Params("postId"))
	if err != nil {
		return err
	}
	var condition = map[string]interface{}{"post_id": postId}
	if err := database.DB.Where(condition).First(&post).Error; err != nil {
		return utils.ErrorResponse(c, utils.GetError)
	}
	var curViews uint = post.Views

	// Update
	if err := database.DB.Model(&post).Update("views", curViews+1).Error; err != nil {
		return utils.ResponseBody(c, utils.UpdateError)
	}

	// Update the view count of this post for current user
	var popularity models.Popularity

	// Search
	condition = map[string]interface{}{"post_id": postId, "user_id": user_id}
	res := database.DB.Where(condition).Limit(1).Find(&popularity)
	if res.RowsAffected == 0 {
		// Create new record as user info has not been recorded yet for this post
		popularity = models.Popularity{
			UserId: user_id,
			PostId: postId,
		}
		if err := database.DB.Create(&popularity).Error; err != nil {
			return utils.ResponseBody(c, utils.CreateError)
		}
	} else {
		// Update view count
		if err := database.DB.Model(&popularity).Update("views", popularity.Views+1).Error; err != nil {
			return utils.ResponseBody(c, utils.UpdateError)
		}
	}

	return utils.ResponseBody(c, utils.Success)
}
