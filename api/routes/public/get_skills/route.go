package get_skills

import (
	"mjrc/core/chix"
	"mjrc/core/postgres"
	"net/http"
)

const (
	Path   = "/skills"
	Method = http.MethodGet
)

func Route(db postgres.DB) *chix.Route {
	hdlr := &handler{db}
	return chix.NewRoute(Path, Method, hdlr.getSkillsForUser)
}
