package app

import (
	"github.com/diyor200/gin-middleware-blogpost/config"
	"github.com/diyor200/gin-middleware-blogpost/internal/controller"
	"github.com/diyor200/gin-middleware-blogpost/internal/middleware"
	"github.com/diyor200/gin-middleware-blogpost/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Run() {
	cfg := config.NewConfig()

	db, err := sqlx.Connect("postgres", cfg.DBUrl)
	//db.AutoMigrate(&entity.User{})
	//db.AutoMigrate(&entity.Blog{})
	if err != nil {
		panic(err)
	}
	repo := repository.NewRepo(db)
	c := controller.NewController(repo)

	router := gin.Default()

	router.GET("/users", c.GetUsers)
	router.GET("/posts", c.GetPosts)
	router.GET("/posts/:post_id", c.GetPost)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", c.SignUp)
		auth.POST("/sign-in", c.SignIn)
	}

	actions := router.Group("/action")
	actions.Use(middleware.CheckUser())
	actions.POST("/create/post", c.CreatePost)
	d := actions.Group("/delete")
	{
		d.POST("/user", c.DeleteUser)
		d.POST("/post/:post_id", c.DeletePost)
	}
	actions.POST("/edit/post", c.EditPost)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
