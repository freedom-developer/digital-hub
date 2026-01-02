package user

import "myapp/middleware"

func (us *UserService) RegisterRouters() {
	userGroup := us.rg

	userGroup.POST("/register", us.RegisterUser)
	userGroup.POST("/login", us.Login)

	authGroup := userGroup.Use(middleware.AuthMiddleware())
	authGroup.GET("/me", us.GetUserProfile)

	// userGroup.GET("/:id", us.GetUser)
	// userGroup.PUT("/:id", us.UpdateUser)
	// userGroup.DELETE("/:id", us.DeleteUser)

}
