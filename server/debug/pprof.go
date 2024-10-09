package debug

import (
	"github.com/chloyka/ginannot"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
)

type DebugControllerInterface struct {
	PPROF ginannot.Route `gin:"GET /debug/pprof/*pprof"`
}

type DebugController struct {
	DebugControllerInterface
}

func NewDebugController() ginannot.Handler {
	return &DebugController{}
}

func (c *DebugController) PPROF(ctx *gin.Context) {
	http.DefaultServeMux.ServeHTTP(ctx.Writer, ctx.Request)
}
