package update_skill

import (
	"mjrc/core/postgres"
	"mjrc/web/chix"
	"net/http"
)

const (
	Path   = "/{id}"
	Method = http.MethodPut
)

func Route(db postgres.DB) *chix.Route {
	hdlr := &handler{db}
	return chix.NewRoute(Path, Method, hdlr.updateSkillForAdmin)
}
