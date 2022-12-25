package database

import (
	"backend/config"
	"backend/models"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(postgres.New(postgres.Config{
		DSN: config.GetPostgresConnectionStr(),
	}), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}
	DB = connection
	ClearTables(connection)
	CreateTables(connection)
	AddDummyData()
}

func CreateTables(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Post{},
		&models.Comment{}, &models.Popularity{})
}

func AddDummyData() {
	AddUser("abby@test.com", "abby123456", "abby", 1, "dark")
	AddUser("bob@test.com", "bob123456", "bob", 1, "red")
	AddUser("cassie@test.com", "cassie123456", "cassie", 2, "pink")
	AddUser("david@test.com", "david123456", "david", 2, "blue")
	AddUser("emily@test.com", "emily123456", "emily", 1, "violet")

	AddPost(1, "How to live healthy?", "Drink more water, get enough sleep, eat more fruits and vegetables and cut down on processed food", "food,lifestyle,health tips")
	AddPost(2, "Sleep deprivation", "Sleep deprivation is when you don’t get the sleep you need, and it is It’s estimated to affect around one-third of American adults , a problem that has only worsened in recent years. Lack of sleep directly affects how we think and feel. While the short-term impacts are more noticeable, chronic sleep deprivation can heighten the long-term risk of physical and mental health problems.", "lifestyle,sleep,health tips")
	AddPost(3, "Worst Foods to Eat for Your Health", "Foods with added sugar or salt or refined carbohydrates, processed meats", "food,groceries")

	AddComment(2, 1, "Maintaining an exercising regime is the most important")
	AddComment(1, 1, "I think so too")
	AddComment(2, 2, "We should have 7 hours of uninterrupted rest")
	AddComment(4, 2, "This can affect our health adversely")
	AddComment(4, 1, "I agree")

	AddPopularity(1, 1, 1, false)
	AddPopularity(2, 1, 1, true)
	AddPopularity(2, 2, 1, false)
	AddPopularity(3, 2, 1, true)
	AddPopularity(4, 2, 2, true)

}

func AddUser(email string, password string, name string, access_type uint, avatar_color string) {
	pw, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	user := models.User{
		Email:       email,
		Password:    pw,
		Name:        name,
		AccessType:  access_type,
		AvatarColor: avatar_color,
	}

	DB.Create(&user)
}

func AddPost(user_id uint, title string, body string, categories string) {
	post := models.Post{
		UserId:     user_id,
		Title:      title,
		Body:       body,
		Categories: categories,
		CreatedDt:  time.Now().Unix(),
	}

	DB.Create(&post)
}

func AddComment(user_id uint, post_id uint, content string) {
	comment := models.Comment{
		UserId:    user_id,
		PostId:    post_id,
		Content:   content,
		CreatedDt: time.Now().Unix(),
	}

	DB.Create(&comment)
}

func AddPopularity(user_id uint, post_id uint, view uint, like bool) {
	popularity := models.Popularity{
		UserId: user_id,
		PostId: post_id,
		Views:  view,
		Likes:  like,
	}

	DB.Create(&popularity)
}

// Remove all records from all tables.
func ClearTables(db *gorm.DB) {
	cols := [4]string{"popularities", "comments", "posts", "users"}
	for _, col := range cols {
		s := fmt.Sprintf("DROP TABLE IF EXISTS %s", col)
		db.Exec(s)
	}
}
