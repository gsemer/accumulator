package application

import (
	"testing"

	"github.com/gorilla/mux"
)

func Test_routes_exist(t *testing.T) {
	testApp := Config{}

	testRoutes := testApp.Routes()
	muxRoutes := testRoutes.(*mux.Router)

	routes := []string{"/add", "/state", "/tsp"}

	for _, r := range routes {
		routeExists(t, muxRoutes, r)
	}
}

func routeExists(t *testing.T, router *mux.Router, route string) {
	t.Helper()

	found := false

	err := router.Walk(func(muxRouter *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		path, err := muxRouter.GetPathTemplate()
		if err != nil {
			return nil
		}

		if path == route {
			found = true
		}
		return nil
	})
	if err != nil {
		t.Fatalf("error walking routes: %v", err)
	}

	if !found {
		t.Errorf("route %q not found", route)
	}
}
