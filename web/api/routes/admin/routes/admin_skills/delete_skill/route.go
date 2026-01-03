package delete_skill

import (
	"mjrc/core/runtime"
	"mjrc/web/chix"
	"net/http"
)

const (
	Path   = "/{id}"
	Method = http.MethodDelete
)

func Route(deps runtime.Dependencies) *chix.Route {
	hdlr := &handler{deps.DB()}
	return chix.NewRoute(Path, Method, hdlr.deleteSkillForAdmin)
}
