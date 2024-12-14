package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/mocks"
	"github.com/stretchr/testify/require"
)

func TestRegistration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockL := mocks.NewMockLogger(ctrl)
	mockUs := mocks.NewMockUserStorager(ctrl)

	type args struct {
		ctx      context.Context
		l        log.Logger
		us       userStorager
		login    string
		password string
	}
	tests := []struct {
		name      string
		args      args
		wantToken string
		wantErr   error
	}{
		{
			name: "Normal user registration",
			args: args{
				ctx:      ctx,
				l:        mockL,
				us:       mockUs,
				login:    "user1",
				password: "password",
			},
			wantToken: "token1",
			wantErr:   nil,
		},
		{
			name: "Registration with empty password",
			args: args{
				ctx:      ctx,
				l:        mockL,
				us:       mockUs,
				login:    "user1",
				password: "",
			},
			wantToken: "",
			wantErr:   ErrUserRegistration,
		},
		{
			name: "Registration with empty login",
			args: args{
				ctx:      ctx,
				l:        mockL,
				us:       mockUs,
				login:    "",
				password: "password",
			},
			wantToken: "",
			wantErr:   ErrUserRegistration,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUs.EXPECT().AddNewUser(tt.args.ctx, tt.args.l, tt.args.login, tt.args.password).Return(tt.wantToken, tt.wantErr)
			mockL.EXPECT().Info("registration new user", "login", tt.args.login, "password", tt.args.password)
			if tt.wantErr != nil {
				mockL.EXPECT().Error("registration user error", "login", tt.args.login, "error", tt.wantErr)
			} else {
				mockL.EXPECT().Info("registration complited", "login", tt.args.login, "token", tt.wantToken)
			}

			gotToken, gotErr := Registration(ctx, tt.args.l, tt.args.us, tt.args.login, tt.args.password)

			require.Equal(t, tt.wantToken, gotToken)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Registration() error = %v, wanted %v", gotErr, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthentication(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockL := mocks.NewMockLogger(ctrl)
	mockUs := mocks.NewMockUserStorager(ctrl)

	type args struct {
		ctx      context.Context
		l        log.Logger
		us       userStorager
		login    string
		password string
	}
	tests := []struct {
		name      string
		args      args
		wantToken string
		wantErr   error
	}{
		{
			name: "Successful authetication",
			args: args{
				ctx:      ctx,
				l:        mockL,
				us:       mockUs,
				login:    "user1",
				password: "password1",
			},
			wantToken: "token1",
			wantErr:   nil,
		},
		{
			name: "Authetication with empty password",
			args: args{
				ctx:      ctx,
				l:        mockL,
				us:       mockUs,
				login:    "user1",
				password: "",
			},
			wantToken: "",
			wantErr:   ErrUserAuthentication,
		},
		{
			name: "Authetication with empty login",
			args: args{
				ctx:      ctx,
				l:        mockL,
				us:       mockUs,
				login:    "",
				password: "password1",
			},
			wantToken: "",
			wantErr:   ErrUserAuthentication,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUs.EXPECT().LoginUser(tt.args.ctx, tt.args.l, tt.args.login, tt.args.password).Return(tt.wantToken, tt.wantErr)
			mockL.EXPECT().Info("authentication user", "login", tt.args.login, "password", tt.args.password)
			if tt.wantErr != nil {
				mockL.EXPECT().Error("authentication user error", "login", tt.args.login, "error", tt.wantErr)
			} else {
				mockL.EXPECT().Info("authentication complited", "login", tt.args.login, "token", tt.wantToken)
			}

			gotToken, gotErr := Authentication(tt.args.ctx, tt.args.l, tt.args.us, tt.args.login, tt.args.password)

			require.Equal(t, tt.wantToken, gotToken)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Authenication() error = %v, wanted %v", gotErr, tt.wantErr)
				return
			}
		})
	}
}
