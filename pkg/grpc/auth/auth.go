package auth

import (
	"context"
	"errors"

	ssov1 "github.com/Pinkman-77/Protobuf/gen/go/sso"
	"github.com/Pinkman-77/records-restapi/pkg/grpc_service/auth"
	"github.com/Pinkman-77/records-restapi/pkg/models"
	"github.com/Pinkman-77/records-restapi/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	ssov1.UnimplementedAuthServer
	auth Auth
	appProvider AppProvider
}

type AppProvider interface {
	CreateApp(ctx context.Context, name, secret string) (int64, error)
	GetApp(ctx context.Context, id int64) (models.App, error)
}
type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	Admin(ctx context.Context, userID int64) (bool, error)
}

func Register(grpc *grpc.Server, auth *auth.Auth ) {
	ssov1.RegisterAuthServer(grpc, &Server{auth: auth})
}

const emptyValue = 0

func (s *Server) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
	}

	func (s *Server) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
		if err := validateRegister(req); err != nil {
			return nil, err
		}
	
		userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
		if err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
	
		return &ssov1.RegisterResponse{UserId: userID}, nil
	
	}

	func (s *Server) Admin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
		if err := validateAdmin(req); err != nil {
			return nil, err
		}
	
		isAdmin, err := s.auth.Admin(ctx, req.GetUserId())
	
		if err != nil {
			if errors.Is(err, storage.ErrAppNotFound) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
	
		return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
	}

	func validateLogin(req *ssov1.LoginRequest) error {
		if req.GetEmail() == "" {
			return status.Error(codes.InvalidArgument, "email is not provided")
		}
	
		if req.GetPassword() == "" {
			return status.Error(codes.InvalidArgument, "password is not provided")
		}
	
		if req.GetAppId() == emptyValue {
			return status.Error(codes.InvalidArgument, "app id is not provided")
		}
	
		return nil
	}

	func validateRegister(req *ssov1.RegisterRequest) error {
		if req.GetEmail() == "" {
			return status.Error(codes.InvalidArgument, "email is not provided")
		}
	
		if req.GetPassword() == "" {
			return status.Error(codes.InvalidArgument, "password is not provided")
		}
	
		return nil
	}

	func validateAdmin(req *ssov1.IsAdminRequest) error {
		if req.GetUserId() == emptyValue {
			return status.Error(codes.InvalidArgument, "user id is not provided")
		}
	
		return nil
	}

	func (s *Server) GetApp(ctx context.Context, req *ssov1.GetAppRequest) (*ssov1.GetAppResponse, error) {
		if err := validateGetApp(req); err != nil {
			return nil, err
		}
	
		app, err := s.appProvider.GetApp(ctx, req.GetId())
	
		if err != nil {
			if errors.Is(err, storage.ErrAppNotFound) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
	
		return &ssov1.GetAppResponse{
			Name: app.Name,
			Secret: app.Secret,
		}, nil
	}

	func validateGetApp(req *ssov1.GetAppRequest) error {
		if req.GetId() == emptyValue {
			return status.Error(codes.InvalidArgument, "app id is not provided")
		}
	
		return nil
	}

	func (s *Server) CreateApp(ctx context.Context, req *ssov1.CreateAppRequest) (*ssov1.CreateAppResponse, error) {
		if err := validateCreateApp(req); err != nil {
			return nil, err
		}
		app, err := s.appProvider.CreateApp(ctx, req.GetName(), req.GetSecret())
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &ssov1.CreateAppResponse{Id: app}, nil
	}

	func validateCreateApp(req *ssov1.CreateAppRequest) error {
		if req.GetName() == "" {
			return status.Error(codes.InvalidArgument, "app name is not provided")
		}
		if req.GetSecret() == "" {
			return status.Error(codes.InvalidArgument, "app secret is not provided")
		}
		return nil
	}