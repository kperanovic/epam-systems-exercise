package storage

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/kperanovic/epam-systems/api/v1/types"
)

func generateCompany(uid uuid.UUID) *types.Company {
	return &types.Company{
		ID:          uid,
		Name:        "test-company",
		Description: "desc",
		Employees:   10,
		Registered:  true,
	}

}

func TestNewMemoryStorage(t *testing.T) {
	tests := []struct {
		name string
		want *memoryStorage
	}{
		{
			name: "Test NewMemoryStorage()",
			want: &memoryStorage{
				store: make(map[uuid.UUID]*types.Company),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryStorage_Connect(t *testing.T) {
	tests := []struct {
		name    string
		m       *memoryStorage
		wantErr bool
	}{
		{
			name:    "Test memory storage Connect()",
			m:       &memoryStorage{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memoryStorage_SaveCompany(t *testing.T) {
	type args struct {
		company *types.Company
	}
	tests := []struct {
		name    string
		m       *memoryStorage
		args    args
		wantErr bool
	}{
		{
			name: "Test memory storage SaveCompany()",
			m: &memoryStorage{
				store: make(map[uuid.UUID]*types.Company),
			},
			args: args{
				company: generateCompany(uuid.New()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SaveCompany(tt.args.company); (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.SaveCompany() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memoryStorage_GetCompany(t *testing.T) {
	uid, err := uuid.NewRandom()
	if err != nil {
		t.Error("error creating uuid")
		return
	}

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		m       *memoryStorage
		args    args
		want    *types.Company
		wantErr bool
	}{
		{
			name: "Test memory storage GetCompany()",
			m: &memoryStorage{
				store: make(map[uuid.UUID]*types.Company),
			},
			args: args{
				id: uid,
			},
			want:    generateCompany(uid),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SaveCompany(generateCompany(tt.args.id))

			got, err := tt.m.GetCompany(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.GetCompany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memoryStorage.GetCompany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryStorage_UpdateCompany(t *testing.T) {
	uid, err := uuid.NewRandom()
	if err != nil {
		t.Error("error creating uuid")
		return
	}

	update := generateCompany(uid)
	update.Employees = 500

	type args struct {
		id      uuid.UUID
		company *types.Company
	}
	tests := []struct {
		name    string
		m       *memoryStorage
		args    args
		wantErr bool
	}{
		{
			name: "Test memory storage UpdateCompany()",
			m: &memoryStorage{
				store: make(map[uuid.UUID]*types.Company),
			},
			args: args{
				id:      uid,
				company: update,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SaveCompany(generateCompany(tt.args.id))

			if err := tt.m.UpdateCompany(tt.args.id, tt.args.company); (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.UpdateCompany() error = %v, wantErr %v", err, tt.wantErr)
			}

			expected := generateCompany(tt.args.id)
			expected.Employees = 500

			got, err := tt.m.GetCompany(tt.args.id)
			if (err != nil) && got != expected {
				t.Errorf("memoryStorage.GetCompany() error = %v, wantErr %v", got, expected)
				return
			}
		})
	}
}

func Test_memoryStorage_DeleteCompany(t *testing.T) {
	uid, err := uuid.NewRandom()
	if err != nil {
		t.Error("error creating uuid")
		return
	}

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		m       *memoryStorage
		args    args
		wantErr bool
	}{
		{
			name: "Test memory store DeleteCompany()",
			m: &memoryStorage{
				store: make(map[uuid.UUID]*types.Company),
			},
			args: args{
				id: uid,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SaveCompany(generateCompany(tt.args.id))

			if err := tt.m.DeleteCompany(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.DeleteCompany() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := tt.m.GetCompany(tt.args.id)
			if (err != nil) && got != nil {
				t.Errorf("memoryStorage.GetCompany() error = %v, wantErr %v", got, nil)
				return
			}
		})
	}
}
