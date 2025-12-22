package delete_skill

import (
	"mjrc/core/chix"
	"mjrc/core/postgres"
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
