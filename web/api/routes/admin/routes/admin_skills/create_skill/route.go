package create_skill

import (
	"mjrc/core/postgres"
	"mjrc/web/chix"
	"net/http"
)

const (
	Path   = "/"
	Method = http.MethodPost
)

func Route(db postgres.DB) *chix.Route {
	hdlr := &handler{db}
	return chix.NewRoute(Path, Method, hdlr.createSkillForAdmin)
}
