package create_skill

import (
	"context"
	"mjrc/core/postgres"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestIntegration_CreateSkill(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(db).Register(r)

	panic("complete implementation")

	//junie
	// complete implementation
}
