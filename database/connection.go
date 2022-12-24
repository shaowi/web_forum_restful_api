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
	AddUser("abby@test.com", "12356", "abby", 1, "dark")
	AddUser("bob@test.com", "123456", "bob", 1, "red")
	AddUser("cassie@test.com", "456789", "cassie", 2, "pink")
	AddUser("shaogamers@gmail.com", "719346", "shao", 1, "grape")
	AddUser("shao.lee@cldcvr.com", "719346", "leeshaowee", 1, "violet")

	AddPost(4, "title1", "this is some body 1", "food,groceries,design")
	AddPost(1, "title2", "this is some body 2", "sports,groceries,design")
	AddPost(2, "title3", "this is some body 3", "shopping,groceries,entertainment")

	AddComment(2, 1, "this is some content 1")
	AddComment(1, 1, "this is some content 2")
	AddComment(2, 2, "this is some content 3")
	AddComment(4, 2, "this is some content 4")
	AddComment(4, 1, "this is some content 5")

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
