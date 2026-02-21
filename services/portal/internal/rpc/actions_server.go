package rpc

import (
	"context"

	"connectrpc.com/connect"

	portalv1 "github.com/jredh-dev/nexus/gen/portal/v1"
	"github.com/jredh-dev/nexus/gen/portal/v1/portalv1connect"
	"github.com/jredh-dev/nexus/services/portal/internal/actions"
	"github.com/jredh-dev/nexus/services/portal/internal/auth"
)

// ActionsServer implements portalv1connect.ActionsServiceHandler.
type ActionsServer struct {
	portalv1connect.UnimplementedActionsServiceHandler

	registry *actions.Registry
	auth     *auth.Service
}

// NewActionsServer creates an ActionsService Connect handler.
func NewActionsServer(registry *actions.Registry, authService *auth.Service) *ActionsServer {
	return &ActionsServer{registry: registry, auth: authService}
}

func (s *ActionsServer) Search(
	ctx context.Context,
	req *connect.Request[portalv1.SearchActionsRequest],
) (*connect.Response[portalv1.SearchActionsResponse], error) {
	query := req.Msg.Query

	// Determine auth context from session cookie (best-effort).
	searchCtx := actions.SearchContext{}
	sessionID := extractSessionCookie(req.Header().Get("Cookie"))
	if sessionID != "" {
		if user, _, err := s.auth.ValidateSession(sessionID); err == nil && user != nil {
			searchCtx.LoggedIn = true
			searchCtx.IsAdmin = user.IsAdmin()
		}
	}

	results := s.registry.Search(query, searchCtx)

	protoActions := make([]*portalv1.Action, len(results))
	for i, a := range results {
		protoActions[i] = &portalv1.Action{
			Id:          a.ID,
			Type:        string(a.Type),
			Title:       a.Title,
			Description: a.Description,
			Target:      a.Target,
		}
	}

	return connect.NewResponse(&portalv1.SearchActionsResponse{
		Actions: protoActions,
	}), nil
}
