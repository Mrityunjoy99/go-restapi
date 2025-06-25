package appserver

import (
	"github.com/Mrityunjoy99/sample-go/src/application"
	"github.com/Mrityunjoy99/sample-go/src/application/admin"
	"github.com/Mrityunjoy99/sample-go/src/application/employee"
	"github.com/Mrityunjoy99/sample-go/src/application/healthcheck"
	"github.com/Mrityunjoy99/sample-go/src/application/user"
	"github.com/Mrityunjoy99/sample-go/src/common/constant"
	"github.com/Mrityunjoy99/sample-go/src/deployment/middleware"
	"github.com/Mrityunjoy99/sample-go/src/domain/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.Engine, appService application.Service, domainService service.ServiceRegistry) {
	globalGroup := g.Group("/")
	globalGroup.Use(middleware.LoggerMiddleware())

	registerHealthCheckRoutes(globalGroup)
	registerUserRoutes(globalGroup, appService, domainService)
	registerAdminRouteGroup(globalGroup, appService, domainService)
	registerEmployeeRouterGroup(globalGroup, appService)
}

func registerAdminRouteGroup(g *gin.RouterGroup, appService application.Service, domainService service.ServiceRegistry) {
	if appService.AdminService == nil {
		panic("AdminService is required for admin routes")
	}

	adminGroup := g.Group("/admin")
	// TODO: Implement and uncomment admin authentication middleware
	// adminGroup.Use(middleware.AdminAuth())
	adminController := admin.NewController(appService.AdminService)
	adminGroup.POST("/generate-token", middleware.AuthMiddleware(domainService.JwtService, constant.UserTypeAdmin), adminController.GenerateToken)
}

func registerHealthCheckRoutes(g *gin.RouterGroup) {
	healthCheckController := healthcheck.NewController()
	g.GET("/health-check", healthCheckController.HealthCheck)
}

func registerUserRoutes(g *gin.RouterGroup, appService application.Service, domainService service.ServiceRegistry) {
	userController := user.NewController(appService.UserService)
	g.GET("/user/:id", middleware.AuthMiddleware(domainService.JwtService, constant.UserTypeUser), userController.GetUserById)
	g.POST("/user", middleware.AuthMiddleware(domainService.JwtService, constant.UserTypeUser), userController.CreateUser)
	g.PUT("/user/:id", middleware.AuthMiddleware(domainService.JwtService, constant.UserTypeUser), userController.UpdateUser)
	g.DELETE("/user/:id", middleware.AuthMiddleware(domainService.JwtService, constant.UserTypeManager), userController.DeleteUser)
}

func registerEmployeeRouterGroup(g *gin.RouterGroup, appService application.Service) {
	employeeController := employee.NewController(appService.EmployeeService)
	g.GET("/employee/:id", employeeController.GetEmployeeById)
}
