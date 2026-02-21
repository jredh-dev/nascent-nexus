package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"connectrpc.com/connect"

	portalv1 "github.com/jredh-dev/nexus/gen/portal/v1"
	"github.com/jredh-dev/nexus/gen/portal/v1/portalv1connect"
	"github.com/jredh-dev/nexus/services/portal/config"
	"github.com/jredh-dev/nexus/services/portal/internal/auth"
	"github.com/jredh-dev/nexus/services/portal/pkg/models"
)

// AuthServer implements portalv1connect.AuthServiceHandler.
type AuthServer struct {
	portalv1connect.UnimplementedAuthServiceHandler

	auth *auth.Service
	cfg  *config.Config
}

// NewAuthServer creates an AuthService Connect handler.
func NewAuthServer(authService *auth.Service, cfg *config.Config) *AuthServer {
	return &AuthServer{auth: authService, cfg: cfg}
}

func (s *AuthServer) Login(
	ctx context.Context,
	req *connect.Request[portalv1.LoginRequest],
) (*connect.Response[portalv1.LoginResponse], error) {
	email := strings.TrimSpace(req.Msg.Email)
	password := req.Msg.Password

	if email == "" || password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email and password are required"))
	}

	// Extract peer info from headers for session tracking.
	ipAddress := req.Header().Get("X-Real-IP")
	if ipAddress == "" {
		ipAddress = req.Header().Get("X-Forwarded-For")
	}
	userAgent := req.Header().Get("User-Agent")

	sessionID, err := s.auth.Login(email, password, ipAddress, userAgent)
	if err != nil {
		log.Printf("Login failed for %s: %v", email, err)
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid email or password"))
	}

	// Set session cookie via response header.
	resp := connect.NewResponse(&portalv1.LoginResponse{
		SessionId: sessionID,
	})
	resp.Header().Set("Set-Cookie", s.sessionCookie(sessionID))

	return resp, nil
}

func (s *AuthServer) Signup(
	ctx context.Context,
	req *connect.Request[portalv1.SignupRequest],
) (*connect.Response[portalv1.SignupResponse], error) {
	username := strings.TrimSpace(req.Msg.Username)
	email := strings.TrimSpace(req.Msg.Email)
	phone := strings.TrimSpace(req.Msg.Phone)
	password := req.Msg.Password
	name := strings.TrimSpace(req.Msg.Name)

	if username == "" || email == "" || phone == "" || password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username, email, phone, and password are required"))
	}

	user, err := s.auth.Signup(username, email, phone, password, name)
	if err != nil {
		log.Printf("Signup failed for %s: %v", email, err)

		switch {
		case errors.Is(err, auth.ErrUsernameTaken):
			return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("username is already taken"))
		case errors.Is(err, auth.ErrEmailTaken):
			return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("an account with this email already exists"))
		case errors.Is(err, auth.ErrPhoneTaken):
			return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("an account with this phone number already exists"))
		default:
			return nil, connect.NewError(connect.CodeInternal, errors.New("signup failed"))
		}
	}

	// Auto-login after signup.
	ipAddress := req.Header().Get("X-Real-IP")
	if ipAddress == "" {
		ipAddress = req.Header().Get("X-Forwarded-For")
	}
	userAgent := req.Header().Get("User-Agent")

	sessionID, err := s.auth.Login(email, password, ipAddress, userAgent)
	if err != nil {
		log.Printf("Auto-login after signup failed for %s: %v", email, err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("account created but auto-login failed"))
	}

	resp := connect.NewResponse(&portalv1.SignupResponse{
		SessionId: sessionID,
		User:      userToProto(user),
	})
	resp.Header().Set("Set-Cookie", s.sessionCookie(sessionID))

	return resp, nil
}

func (s *AuthServer) Logout(
	ctx context.Context,
	req *connect.Request[portalv1.LogoutRequest],
) (*connect.Response[portalv1.LogoutResponse], error) {
	sessionID := extractSessionCookie(req.Header().Get("Cookie"))
	if sessionID != "" {
		_ = s.auth.Logout(sessionID)
	}

	resp := connect.NewResponse(&portalv1.LogoutResponse{})
	// Clear the cookie.
	resp.Header().Set("Set-Cookie", "session=; Path=/; Max-Age=-1; HttpOnly; SameSite=Lax")

	return resp, nil
}

func (s *AuthServer) GetSession(
	ctx context.Context,
	req *connect.Request[portalv1.GetSessionRequest],
) (*connect.Response[portalv1.GetSessionResponse], error) {
	sessionID := extractSessionCookie(req.Header().Get("Cookie"))
	if sessionID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no session"))
	}

	user, _, err := s.auth.ValidateSession(sessionID)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid session"))
	}

	sessions, err := s.auth.GetSessionsByUserID(user.ID)
	if err != nil {
		log.Printf("Error fetching sessions for user %s: %v", user.ID, err)
		sessions = nil
	}

	return connect.NewResponse(&portalv1.GetSessionResponse{
		User:     userToProto(user),
		Sessions: sessionsToProto(sessions),
	}), nil
}

func (s *AuthServer) MagicLogin(
	ctx context.Context,
	req *connect.Request[portalv1.MagicLoginRequest],
) (*connect.Response[portalv1.MagicLoginResponse], error) {
	token := strings.TrimSpace(req.Msg.Token)
	if token == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("token is required"))
	}

	ipAddress := req.Header().Get("X-Real-IP")
	if ipAddress == "" {
		ipAddress = req.Header().Get("X-Forwarded-For")
	}
	userAgent := req.Header().Get("User-Agent")

	sessionID, err := s.auth.ValidateMagicToken(token, ipAddress, userAgent)
	if err != nil {
		log.Printf("Magic login failed: %v", err)
		if errors.Is(err, auth.ErrInvalidMagicToken) {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid or expired magic login link"))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.New("magic login failed"))
	}

	resp := connect.NewResponse(&portalv1.MagicLoginResponse{
		SessionId: sessionID,
	})
	resp.Header().Set("Set-Cookie", s.sessionCookie(sessionID))

	return resp, nil
}

func (s *AuthServer) GenerateMagicLink(
	ctx context.Context,
	req *connect.Request[portalv1.GenerateMagicLinkRequest],
) (*connect.Response[portalv1.GenerateMagicLinkResponse], error) {
	// Require admin session.
	sessionID := extractSessionCookie(req.Header().Get("Cookie"))
	if sessionID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	user, _, err := s.auth.ValidateSession(sessionID)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid session"))
	}
	if !user.IsAdmin() {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("admin access required"))
	}

	email := strings.TrimSpace(req.Msg.Email)
	if email == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email is required"))
	}

	token, err := s.auth.CreateMagicToken(email)
	if err != nil {
		log.Printf("Failed to generate magic link for %s: %v", email, err)
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate magic link"))
	}

	// Build the magic link URL using the Host header.
	scheme := "https"
	if s.cfg.Server.Env != "production" {
		scheme = "http"
	}
	host := req.Header().Get("Host")
	link := fmt.Sprintf("%s://%s/auth/magic?token=%s", scheme, host, token)

	return connect.NewResponse(&portalv1.GenerateMagicLinkResponse{
		MagicLink: link,
	}), nil
}

// --- helpers ---

func (s *AuthServer) sessionCookie(sessionID string) string {
	secure := ""
	if s.cfg.Server.Env == "production" {
		secure = "; Secure"
	}
	return fmt.Sprintf("session=%s; Path=/; Max-Age=%d; HttpOnly; SameSite=Lax%s",
		sessionID, s.cfg.Session.MaxAge, secure)
}

func extractSessionCookie(cookieHeader string) string {
	for _, part := range strings.Split(cookieHeader, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "session=") {
			return strings.TrimPrefix(part, "session=")
		}
	}
	return ""
}

func userToProto(u *models.User) *portalv1.User {
	if u == nil {
		return nil
	}
	return &portalv1.User{
		Id:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.PhoneNumber,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LastLogin: u.LastLoginAt.Format("2006-01-02T15:04:05Z"),
	}
}

func sessionsToProto(sessions []models.Session) []*portalv1.Session {
	if sessions == nil {
		return nil
	}
	out := make([]*portalv1.Session, len(sessions))
	for i, s := range sessions {
		out[i] = &portalv1.Session{
			Id:        s.ID,
			IpAddress: s.IPAddress,
			UserAgent: s.UserAgent,
			CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z"),
			ExpiresAt: s.ExpiresAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return out
}
