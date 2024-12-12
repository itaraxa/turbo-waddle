package services

import (
	"context"
	"errors"
	"testing"

	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegistration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := mocks.NewMockLogger(ctrl)
	m := mocks.NewMockUserStorager(ctrl)

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
				l:        l,
				us:       m,
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
				l:        l,
				us:       m,
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
				l:        l,
				us:       m,
				login:    "",
				password: "password",
			},
			wantToken: "",
			wantErr:   ErrUserRegistration,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().AddNewUser(ctx, l, tt.args.login, tt.args.password).Return(tt.wantToken, tt.wantErr)
			l.EXPECT().Info("registration new user", "login", tt.args.login, "password", tt.args.password)
			if tt.wantErr != nil {
				l.EXPECT().Error("registration user error", "login", tt.args.login, "error", tt.wantErr)
			} else {
				l.EXPECT().Info("registration complited", "login", tt.args.login, "token", tt.wantToken)
			}

			gotToken, gotErr := Registration(ctx, l, m, tt.args.login, tt.args.password)

			require.Equal(t, gotToken, tt.wantToken)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Registration() error = %v, wanted %v", gotErr, tt.wantErr)
				return
			}

		})
	}
}
