package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID   string `gorm:"primary_key"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0"`
}

var (
	db         *gorm.DB
	validate   *validator.Validate
)

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	db.AutoMigrate(&User{})
	validate = validator.New()
}

func main() {
	r := gin.Default()
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", GetUsers)
		userRoutes.POST("/", CreateUser)
		userRoutes.PUT("/:id", EditUser)
		userRoutes.DELETE("/:id", DeleteUser)
	}

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err.Error())
	}
}

func GetUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error retrieving users"})
		return
	}
	c.JSON(200, users)
}

func CreateUser(c *gin.Context) {
	var reqBody User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := validate.Struct(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	reqBody.ID = uuid.New().String()

	if err := db.Create(&reqBody).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(201, gin.H{"error": false, "id": reqBody.ID})
}

func EditUser(c *gin.Context) {
	id := c.Param("id")
	var reqBody User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := validate.Struct(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	user.Name = reqBody.Name
	user.Age = reqBody.Age
	if err := db.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error updating user"})
		return
	}

	c.JSON(200, gin.H{"error": false})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&User{}, "id = ?", id).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, gin.H{"error": false})
}
