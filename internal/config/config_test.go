package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestGopherMartConfig_Config(t *testing.T) {
	type fields struct {
		Endpoint             string
		DSN                  string
		AccrualSystemAddress string
		LogLevel             string
		ShowVersion          bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Test correct configuration",
			wantErr: false,
			fields: fields{
				Endpoint:             `http://localhost:8080`,
				DSN:                  `postgres://gophermart:\!qaz2wsx@localhost:5432/gophermart`,
				AccrualSystemAddress: `http://localhost:8081`,
				LogLevel:             LOG_INFO,
				ShowVersion:          false,
			},
		},
		{
			name:    "Test empty endpoint",
			wantErr: true,
			fields: fields{
				Endpoint:             ``,
				DSN:                  `postgres://gophermart:\!qaz2wsx@localhost:5432/gophermart`,
				AccrualSystemAddress: `http://localhost:8081`,
				LogLevel:             LOG_INFO,
				ShowVersion:          false,
			},
		},
		{
			name:    "Test unknown log level",
			wantErr: true,
			fields: fields{
				Endpoint:             `http://localhost:8080`,
				DSN:                  `postgres://gophermart:\!qaz2wsx@localhost:5432/gophermart`,
				AccrualSystemAddress: `http://localhost:8081`,
				LogLevel:             `UNKNOWN`,
				ShowVersion:          false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_via_Flags", tt.name), func(t *testing.T) {
			// flag test
			getGmc := NewGopherMartConfig()
			wantGmc := &GopherMartConfig{
				Endpoint:             tt.fields.Endpoint,
				DSN:                  tt.fields.DSN,
				AccrualSystemAddress: tt.fields.AccrualSystemAddress,
				LogLevel:             tt.fields.LogLevel,
				ShowVersion:          tt.fields.ShowVersion,
			}

			originalArs := os.Args
			defer func() {
				os.Args = originalArs
			}()

			originalCommandLine := flag.CommandLine
			defer func() {
				flag.CommandLine = originalCommandLine
			}()
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			os.Args = []string{
				os.Args[0],
				fmt.Sprintf("-a=%s", tt.fields.Endpoint),
				fmt.Sprintf("-d=%s", tt.fields.DSN),
				fmt.Sprintf("-r=%s", tt.fields.AccrualSystemAddress),
				fmt.Sprintf("-l=%s", tt.fields.LogLevel),
			}

			err := getGmc.Config()
			if !reflect.DeepEqual(getGmc, wantGmc) {
				t.Errorf("Configuration by flags: Wanted config = %v Getted config = %v", wantGmc, getGmc)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GopherMartConfig.Config() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_via_ENV", tt.name), func(t *testing.T) {
			// env test
			getGmc := NewGopherMartConfig()
			wantGmc := &GopherMartConfig{
				Endpoint:             tt.fields.Endpoint,
				DSN:                  tt.fields.DSN,
				AccrualSystemAddress: tt.fields.AccrualSystemAddress,
				LogLevel:             tt.fields.LogLevel,
			}

			origEndpoint, existEndpoint := os.LookupEnv(`RUN_ADDRESS`)
			origDSN, existDSN := os.LookupEnv(`DATABASE_URI`)
			origAccrualSystemAddres, existAccrualSystemAddres := os.LookupEnv(`ACCRUAL_SYSTEM_ADDRESS`)
			defer func() {
				if existEndpoint {
					_ = os.Setenv(`RUN_ADDRESS`, origEndpoint)
				} else {
					_ = os.Unsetenv(`RUN_ADDRESS`)
				}
				if existDSN {
					_ = os.Setenv(`DATABASE_URI`, origDSN)
				} else {
					_ = os.Unsetenv(`DATABASE_URI`)
				}
				if existAccrualSystemAddres {
					_ = os.Setenv(`ACCRUAL_SYSTEM_ADDRESS`, origAccrualSystemAddres)
				} else {
					_ = os.Unsetenv(`ACCRUAL_SYSTEM_ADDRESS`)
				}
			}()

			originalCommandLine := flag.CommandLine
			defer func() {
				flag.CommandLine = originalCommandLine
			}()
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			os.Args = []string{
				os.Args[0],
				fmt.Sprintf("-l=%s", tt.fields.LogLevel),
			}

			os.Setenv(`RUN_ADDRESS`, tt.fields.Endpoint)
			os.Setenv(`DATABASE_URI`, tt.fields.DSN)
			os.Setenv(`ACCRUAL_SYSTEM_ADDRESS`, tt.fields.AccrualSystemAddress)

			err := getGmc.Config()
			if !reflect.DeepEqual(getGmc, wantGmc) {
				t.Errorf("Configuration by env: Wanted config = %v Getted config = %v", wantGmc, getGmc)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GopherMartConfig.Config() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
