package user

func (us *UserService) RegisterRouters() {
	userGroup := us.rg.Group("/users")
	{
		userGroup.POST("/register", us.RegisterUser)
		// userGroup.POST("/login", us.Login)
		// userGroup.GET("/:id", us.GetUser)
		// userGroup.PUT("/:id", us.UpdateUser)
		// userGroup.DELETE("/:id", us.DeleteUser)
	}
}
