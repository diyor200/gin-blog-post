package controller

import (
	"errors"
	"fmt"
	"github.com/diyor200/gin-middleware-blogpost/internal/entity"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (c *Controller) CreatePost(ctx *gin.Context) {
	var err error
	userId, ok := ctx.Get("user_id")
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, ErrUnauthorized)
		return
	}
	body := entity.BlogInput{}

	if err = ctx.BindJSON(&body); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	postBody := entity.CreateBlogInput{PostBody: body.PostBody, PostTittle: body.PostTittle, UserID: userId.(int)}
	err = c.r.CreatePost(postBody)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Created post"})
	return
}

func (c *Controller) GetPosts(ctx *gin.Context) {
	posts, err := c.r.GetPosts()
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(200, gin.H{"posts": posts})
	return
}

func (c *Controller) GetPost(ctx *gin.Context) {
	postId, err := strconv.Atoi(ctx.Param("post_id"))
	fmt.Println("GetPost - post_id = ", postId)
	if err != nil {
		errorResponse(ctx, 400, err)
		return
	}
	post, err := c.r.GetPost(postId)
	if err != nil || len(post) == 0 {
		errorResponse(ctx, http.StatusBadRequest, errors.New("not found"))
		return
	}

	ctx.JSON(200, gin.H{"blog": post})
	return
}

func (c *Controller) EditPost(ctx *gin.Context) {
	var input entity.BlogInput
	var err error
	userId, ok := ctx.Get("user_id")
	fmt.Println("EditPost userId = ", userId)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, ErrUnauthorized)
		return
	}
	if err = ctx.BindJSON(&input); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err = c.r.EditPost(input, userId.(int))
	log.Println("err = c.r.EditPost(input, userId.(int)) = ", err)
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Updated post"})
	return
}

func (c *Controller) DeletePost(ctx *gin.Context) {
	postId, err := strconv.Atoi(ctx.Param("post_id"))
	userId, ok := ctx.Get("user_id")
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	err = c.r.DeletePost(postId, userId.(int))
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, err)
		return
	}

	ctx.JSON(200, gin.H{"message": "post deleted"})
}

func errorResponse(ctx *gin.Context, status int, err error) {
	ctx.AbortWithStatusJSON(status, map[string]string{"error": fmt.Sprintf("%v", err)})
}
