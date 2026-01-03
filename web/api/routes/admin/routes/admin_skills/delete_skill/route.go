package delete_skill

import (
	"mjrc/core/postgres"
	"mjrc/web/chix"
	"net/http"
)

const (
	Path   = "/{id}"
	Method = http.MethodDelete
)

func Route(db postgres.DB) *chix.Route {
	hdlr := &handler{db}
	return chix.NewRoute(Path, Method, hdlr.deleteSkillForAdmin)
}
