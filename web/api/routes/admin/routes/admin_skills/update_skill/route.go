package update_skill

import (
	"mjrc/core/runtime"
	"mjrc/web/chix"
	"net/http"
)

const (
	Path   = "/{id}"
	Method = http.MethodPut
)

func Route(deps runtime.Dependencies) *chix.Route {
	hdlr := &handler{deps.DB()}
	return chix.NewRoute(Path, Method, hdlr.updateSkillForAdmin)
}
