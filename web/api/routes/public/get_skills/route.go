package get_skills

import (
	"mjrc/core/runtime"
	"mjrc/web/chix"
	"net/http"
)

const (
	Path   = "/skills"
	Method = http.MethodGet
)

func Route(deps runtime.Dependencies) *chix.Route {
	hdlr := &handler{deps.DB()}
	return chix.NewRoute(Path, Method, hdlr.getSkillsForUser)
}
