package routes

import (
	"clubhub-hotel-management/cmd/server/handler"
	"clubhub-hotel-management/internal/franquicia"
	"os"

	_ "clubhub-hotel-management/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router interface {
	MapRoutes()
}

type router struct {
	r       *gin.Engine
	rg      *gin.RouterGroup
	mongodb *mongo.Client
}

func NewRouter(r *gin.Engine, mongoDB *mongo.Client) Router {
	return &router{r: r, mongodb: mongoDB}
}

func (r *router) MapRoutes() {
	r.setGroup()
	r.buildRoutes()
}

func (r *router) setGroup() {
	r.r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.rg = r.r.Group("/api/hotelmagnament/v1")
}

func (r *router) buildRoutes() {
	repository := franquicia.NewRepository(r.mongodb.Database(os.Getenv("MONGODB_DATABASE_NAME")).Collection("franchises"))
	service := franquicia.NewService(repository)
	fHandler := handler.NewUser(service)
	franchises := r.rg.Group("/franchises")
	franchises.POST("/new", fHandler.Create())
	franchises.GET("/all", fHandler.GetAllFranquicias())
	franchises.PUT("/:id", fHandler.UpdateFranquicia())
	franchises.GET("/one/:id", fHandler.GetFranquiciaByID())
	franchises.GET("/location", fHandler.GetByLocation())
	franchises.GET("/daterange", fHandler.GetFranquiciasByDateRange())
	franchises.GET("/name", fHandler.GetFranquiciasByName())
}
