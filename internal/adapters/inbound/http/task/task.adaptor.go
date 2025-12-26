package poll

import (
	"github.com/gin-gonic/gin"
	service "github.com/thealiakbari/task-pool-system/internal/application/task"
)

type Adaptor struct {
	service.TaskHttpApp
}

func (a Adaptor) RegisterRoutes(r *gin.RouterGroup) {
	apiTask := r.Group("/tasks")

	apiTask.POST("", a.MakeCreate())
	apiTask.PUT("/:id", a.MakeUpdate())

	apiTask.GET("/:id", a.MakeGetById())

	apiTask.DELETE("/:id", a.MakeDelete())
	apiTask.DELETE("/purge/:id", a.MakePurge())
}
