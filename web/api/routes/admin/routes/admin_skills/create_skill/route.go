package create_skill

import (
	"mjrc/core/runtime"
	"mjrc/web/chix"
	"net/http"
)

const (
	Path   = "/"
	Method = http.MethodPost
)

func Route(deps runtime.Dependencies) *chix.Route {
	hdlr := &handler{deps.DB()}
	return chix.NewRoute(Path, Method, hdlr.createSkillForAdmin)
}
