package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis server local address
		Password: "",
		DB:       0,
	})
}

func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	val, err := rdb.Get(c, strconv.Itoa(id)).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, gin.H{"data": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	user := User{
		ID:   id,
		Name: val,
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func setUser(c *gin.Context) {
	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	id := strconv.Itoa(user.ID)
	if rdb.Exists(c, id).Val() == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID already exists"})
		return
	}

	err = rdb.Set(c, strconv.Itoa(user.ID), user.Name, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Set new user data successfully"})
}

func deleteUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	key := strconv.Itoa(id)

	if rdb.Exists(c, key).Val() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
		return
	}

	err = rdb.Del(c, key).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User has been deleted successfully"})
}

func main() {
	initRedis()

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", getUserByID)
		v1.POST("/users", setUser)
		v1.DELETE("/users", deleteUserByID)
	}
	r.Run(":8080")
}
