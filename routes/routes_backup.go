package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/luongquochai/promotional-campaign-system/controllers"
// 	"github.com/luongquochai/promotional-campaign-system/middleware"
// )

// func SetupRoutes(router *gin.Engine) {
// 	// User routes
// 	user := router.Group("/user")
// 	{
// 		user.POST("/register", controllers.Register)
// 		user.POST("/login", controllers.Login)
// 	}

// 	// Campaign routes
// 	campaign := router.Group("/campaigns")
// 	campaign.Use(middleware.AuthMiddleware())
// 	{
// 		campaign.POST("", controllers.CreateCampaign)       // Create a new campaign
// 		campaign.GET("", controllers.ListCampaigns)         // List all campaigns base on userID
// 		campaign.GET("/:id", controllers.GetCampaign)       // Get campaign details by ID
// 		campaign.PUT("/:id", controllers.UpdateCampaign)    // Update campaign information
// 		campaign.DELETE("/:id", controllers.DeleteCampaign) // Delete campaign
// 	}
// 	voucher := router.Group("/voucher")
// 	voucher.Use(middleware.AuthMiddleware()) // Apply middleware
// 	{
// 		voucher.POST("/generate", controllers.GenerateVoucher)
// 		voucher.POST("/validate", controllers.ValidateVoucher)
// 	}

// 	purchase := router.Group("/purchase")
// 	purchase.Use(middleware.AuthMiddleware()) // Apply the authentication middleware
// 	{
// 		purchase.POST("/create", controllers.CreatePurchase) // create purchase and wait for confirm -> export bill
// 		purchase.GET("/history", controllers.GetPurchaseHistory)
// 	}
// }
