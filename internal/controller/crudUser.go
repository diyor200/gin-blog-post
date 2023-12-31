package controller

import (
	"errors"
	"github.com/diyor200/gin-middleware-blogpost/internal/entity"
	"github.com/diyor200/gin-middleware-blogpost/pkg/hash"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (c *Controller) SignUp(ctx *gin.Context) {
	input := entity.UserInput{}
	if err := ctx.BindJSON(&input); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	input.Password = hash.Hash(input.Password)
	err := c.r.CreateUser(input)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, errors.New("user already exists"))
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "successfully created"})
	return
}

func (c *Controller) SignIn(ctx *gin.Context) {
	var input entity.SignInInput
	if err := ctx.BindJSON(&input); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := c.r.GetUserByEmail(input.Email)
	log.Println("Signin - controller - c.r.GetUserByEmail - error: ", err)
	log.Println("Signin - controller - c.r.GetUserByEmail - user: ", user)
	if err != nil {
		errorResponse(ctx, 400, err)
		return
	}
	if user.Username == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}
	input.Password = hash.Hash(input.Password)
	if user.Password != input.Password {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bcrypt.CompareHashAndPassword: invalid email or password"})
		return
	}

	tokenString, err := GenerateToken(user.Id)
	if err != nil {
		errorResponse(ctx, 500, err)
		return
	}
	ctx.JSON(200, gin.H{"token": tokenString})
	return
}

func (c *Controller) GetUsers(ctx *gin.Context) {
	var users []entity.User
	users, err := c.r.GetUsers()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(200, &users)
	return
}

func (c *Controller) DeleteUser(ctx *gin.Context) {
	value := ctx.Param("id")
	id, err := strconv.Atoi(value)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = c.r.DeleteUser(id)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(200, gin.H{"message": "deleted"})
	return
}
