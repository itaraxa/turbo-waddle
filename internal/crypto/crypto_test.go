package crypto

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func Test_GenerateSalt(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "Generating 16-bytes salt",
			size:    16,
			wantErr: false,
		},
		{
			name:    "Incorrect syze salt",
			size:    0,
			wantErr: true,
		},
		{
			name:    "Incorrect syze salt",
			size:    9999,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSalt(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSalt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.size && err == nil {
				t.Errorf("GenerateSalt() = %v, want [16]byte", got)
			}
		})
	}
}

func TestGeneratePasswordWithSaltHash(t *testing.T) {
	type args struct {
		salt     []byte
		password []byte
	}
	tests := []struct {
		name     string
		args     args
		wantHash [32]byte
		wantErr  bool
	}{
		{
			name: "Test generation 1",
			args: args{
				salt:     []byte("saltsalt"),
				password: []byte("password"),
			},
			wantHash: convertHashStrToByteArray("1e4b70574b8ec9633ef4a1fc1113fbf5c04b6ec810798554bbc6bfe85beabb4c"),
			wantErr:  false,
		},
		{
			name: "Test generation 2",
			args: args{
				salt:     []byte("b98866d2be6be"),
				password: []byte("744e3866f13c0"),
			},
			wantHash: convertHashStrToByteArray("b175dc6467e38d5e4151878d699adcc6ac033ed5134c7cbe539e9e0f1fdb264c"),
			wantErr:  false,
		},
		{
			name: "Test short salt",
			args: args{
				salt:     []byte("14"),
				password: []byte("744e3866f13c0"),
			},
			wantErr: true,
		},
		{
			name: "Test empty password",
			args: args{
				salt:     []byte("b98866d2be6be"),
				password: []byte(""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePasswordWithSaltHash(tt.args.salt, tt.args.password)
			if !reflect.DeepEqual(got, tt.wantHash) {
				t.Errorf("GeneratePasswordWithSaltHash() = %s, want %s", hex.EncodeToString(got[:]), hex.EncodeToString(tt.wantHash[:]))
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePasswordWithSaltHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func BenchmarkGeneratePasswordWithSaltHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		salt, _ := GenerateSalt(32)
		username := []byte(fmt.Sprintf("user%d", i))
		_, _ = GeneratePasswordWithSaltHash(salt, username)
	}
}

func convertHashStrToByteArray(str string) (hash [32]byte) {
	t, _ := hex.DecodeString(str)
	copy(hash[:], t)
	return
}

func TestCheckPassword(t *testing.T) {
	type args struct {
		salt       []byte
		password   []byte
		storedHash [32]byte
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
		wantErr    bool
	}{
		{
			name: "Valid password",
			args: args{
				salt:       []byte("saltsalt"),
				password:   []byte("password"),
				storedHash: convertHashStrToByteArray("1e4b70574b8ec9633ef4a1fc1113fbf5c04b6ec810798554bbc6bfe85beabb4c"),
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "Invalid password",
			args: args{
				salt:       []byte("saltsalt"),
				password:   []byte("bla-bla-bla"),
				storedHash: convertHashStrToByteArray("1e4b70574b8ec9633ef4a1fc1113fbf5c04b6ec810798554bbc6bfe85beabb4c"),
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name: "Empty password",
			args: args{
				salt:       []byte("saltsalt"),
				password:   []byte(""),
				storedHash: convertHashStrToByteArray("1e4b70574b8ec9633ef4a1fc1113fbf5c04b6ec810798554bbc6bfe85beabb4c"),
			},
			wantResult: false,
			wantErr:    true,
		},
		{
			name: "Incorrect salt",
			args: args{
				salt:       []byte("sl"),
				password:   []byte("password"),
				storedHash: convertHashStrToByteArray("1e4b70574b8ec9633ef4a1fc1113fbf5c04b6ec810798554bbc6bfe85beabb4c"),
			},
			wantResult: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := CheckPassword(tt.args.salt, tt.args.password, tt.args.storedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("CheckPassword() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestGenerateToken64(t *testing.T) {
	tokenStorage := map[string]struct{}{}

	t.Run("Uniq token test", func(t *testing.T) {
		for i := 0; i < 1e6; i++ {
			gotToken, err := GenerateToken64()
			if err != nil {
				t.Errorf("GenerateToken64() error = %v, wantErr %v", err, false)
				return
			}
			if _, ok := tokenStorage[gotToken]; ok {
				t.Error("duplicate token generated")
				return
			}
			tokenStorage[gotToken] = struct{}{}
		}
	})
}

func BenchmarkGenerateToken64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenerateToken64()
		if err != nil {
			b.Errorf("GenerateToken64() error = %v", err)
		}
	}
}

func TestCreateJWT(t *testing.T) {
	type args struct {
		login     string
		secretKey []byte
	}
	tests := []struct {
		name            string
		args            args
		wantTokenString string
		wantErr         bool
	}{
		{
			name: "Generate JWT test: valid data",
			args: args{
				login:     "user1",
				secretKey: []byte("secretKey"),
			},
			wantTokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6InVzZXIxIn0.R7AZOJAL7ouhlfE6K3LFmUPIAsjBxLvzF-Vngmr0QeM",
			wantErr:         false,
		},
		{
			name: "Generate JWT test: empty secret key",
			args: args{
				login:     "user1",
				secretKey: []byte(""),
			},
			wantTokenString: "",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokenString, err := CreateJWT(tt.args.login, tt.args.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTokenString != tt.wantTokenString {
				t.Errorf("CreateJWT() = %v, want %v", gotTokenString, tt.wantTokenString)
			}
		})
	}
}

func TestVerifyJWT(t *testing.T) {
	type args struct {
		tokenString string
		secretKey   []byte
	}
	tests := []struct {
		name      string
		args      args
		wantValid bool
		wantErr   bool
	}{
		{
			name: "Verify JWT test: valid token",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6InVzZXIxIn0.R7AZOJAL7ouhlfE6K3LFmUPIAsjBxLvzF-Vngmr0QeM",
				secretKey:   []byte("secretKey"),
			},
			wantValid: true,
			wantErr:   false,
		},
		{
			name: "Verify JWT test: invalid token",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6InVzZXIxIn0.R7AZOJAL7ouhlfE6K3LFmUPIAsjBxLvzF-Vngmr0QPq",
				secretKey:   []byte("secretKey"),
			},
			wantValid: false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, err := VerifyJWT(tt.args.tokenString, tt.args.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotValid != tt.wantValid {
				t.Errorf("VerifyJWT() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}
