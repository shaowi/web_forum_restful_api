package utils

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

// Wrapper around strconv.ParseUint. Base is set to 10; bitSize is set to 64.
func ParseUint(s string) (uint, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	return uint(num), err
}

func GetJwtToken(cookie string, SecretKey string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
}

func GetCurrentUser(c *fiber.Ctx, SecretKey string) (models.User, error) {
	cookie := c.Cookies("jwt")

	token, err := GetJwtToken(cookie, SecretKey)

	var user models.User

	if err != nil {
		return user, err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	database.DB.Where("user_id = ?", claims.Issuer).First(&user)

	return user, nil
}

func GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var password strings.Builder
	var (
		lowerCharSet   = "abcdedfghijklmnopqrst"
		upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialCharSet = "!@#$%&*"
		numberSet      = "0123456789"
		allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
	)

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func SendEmail(newPw string, emailTo string) {
	from := config.GoDotEnvVariable("EmailFrom")
	password := config.GoDotEnvVariable("EmailFromSMTPPw")

	toEmailAddress := emailTo
	to := []string{toEmailAddress}

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	subject := "Subject: Reset your password\n"
	body := fmt.Sprintf("Your new password is %s.\nPlease change your password.", newPw)
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		panic(err)
	}

}

func GetHistoryPosts(c *fiber.Ctx, condition map[string]interface{}, joinCondition string) error {
	posts := []models.Post{}
	if err := database.DB.Joins(joinCondition).Joins("User").Where(condition).Find(&posts).Error; err != nil {
		return ErrorResponse(c, GetError)
	}

	for i, post := range posts {
		posts[i] = GetPostStats(post)
	}
	return GetRequestResponse(c, posts)
}

func GetPostStats(post models.Post) models.Post {
	postId := post.PostId
	var views uint = post.Views
	var likes uint = post.Likes
	var comments uint = post.Comments
	viewQuery := "SELECT SUM(views) FROM popularities WHERE post_id = ?"
	likeQuery := "SELECT COUNT(likes) FROM popularities WHERE post_id = ? AND likes = true"
	commentQuery := "SELECT COUNT(comment_id) FROM comments WHERE post_id = ?"

	// Check view records exist in popularities
	var popularity models.Popularity
	condition := map[string]interface{}{"post_id": postId}
	viewsRes := database.DB.Where(condition).Limit(1).Find(&popularity)
	if viewsRes.RowsAffected > 0 {
		database.DB.Raw(viewQuery, postId).Scan(&views)
	}
	database.DB.Raw(likeQuery, postId).Scan(&likes)
	database.DB.Raw(commentQuery, postId).Scan(&comments)

	post.Views = views
	post.Likes = likes
	post.Comments = comments

	database.DB.Save(&post)
	return post
}
