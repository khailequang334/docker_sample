package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/khailequang334/docker_sample/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var rdb *redis.Client

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success string `json:"success"`
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis server local address
		Password: "",
		DB:       0,
	})
}

// @Summary 	Get user by ID
// @Description Get user information by ID
// @Produce 	json
// @Param 		id query int true "User ID"
// @Success 	200 {object} User
// @Failure 	400 {object} ErrorResponse
// @Failure 	404 {object} ErrorResponse
// @Failure 	500 {object} ErrorResponse
// @Router 		/users [get]
func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	val, err := rdb.Get(c, strconv.Itoa(id)).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, User{
		ID:   id,
		Name: val,
	})
}

// @Summary 	Create user
// @Description Create new user information
// @Accept 		json
// @Produce 	json
// @Param 		user body User true "User object"
// @Success 	200 {object} SuccessResponse
// @Failure 	400 {object} ErrorResponse
// @Failure 	500 {object} ErrorResponse
// @Router 		/users [post]
func createUser(c *gin.Context) {
	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
		return
	}

	id := strconv.Itoa(user.ID)
	if rdb.Exists(c, id).Val() == 1 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ID already exists"})
		return
	}

	err = rdb.Set(c, strconv.Itoa(user.ID), user.Name, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Success: "Set new user data successfully"})
}

// @Summary 	Delete user by ID
// @Description Delete user by ID
// @Produce 	json
// @Param 		id query int true "User ID"
// @Success 	200 {object} SuccessResponse
// @Failure 	400 {object} ErrorResponse
// @Failure 	404 {object} ErrorResponse
// @Failure 	500 {object} ErrorResponse
// @Router 		/users [delete]
func deleteUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	key := strconv.Itoa(id)

	if rdb.Exists(c, key).Val() == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "ID not found"})
		return
	}

	err = rdb.Del(c, key).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Success: "User has been deleted successfully"})
}

// @title           Gin Web Users Management
// @version         1.0
// @description     A simple web service for testing.
// @host			localhost:8080
// @BasePath  		/v1/
func main() {
	initRedis()

	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/users", getUserByID)
		v1.POST("/users", createUser)
		v1.DELETE("/users", deleteUserByID)
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.Run(":8080")
}
