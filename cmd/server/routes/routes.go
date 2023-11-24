package routes

import (
	"clubhub-hotel-management/cmd/server/handler"
	"clubhub-hotel-management/internal/franquicia"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/basic/docs"
)

type Router interface {
	MapRoutes()
}

type router struct {
	r       *gin.Engine
	rg      *gin.RouterGroup
	mysqldb *sql.DB
}

func NewRouter(r *gin.Engine, mysql *sql.DB) Router {
	return &router{r: r, mysqldb: mysql}
}

func (r *router) MapRoutes() {
	r.setGroup()
	r.buildRoutes()
}

func (r *router) setGroup() {
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.rg = r.r.Group("/api/hotelmagnament/v1")
}

func (r *router) buildRoutes() {

	repository := franquicia.NewRepository(r.mysqldb)
	service := franquicia.NewService(repository)
	fHandler := handler.NewUser(service)
	franchises := r.rg.Group("/franchises")
	franchises.GET("/one")
	franchises.POST("/new", fHandler.Create())
	franchises.GET("/all", fHandler.GetAllFranquicias())
	//r.rg.DELETE("/user/:id", fHandler)
	//r.rg.PUT("/user/:id", fHandler)
}
