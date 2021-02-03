package business

import (
	"reflect"
	"testing"

	"github.com/oguzhn/PlayerRankingService/database"
	"github.com/oguzhn/PlayerRankingService/models"
)

type mockDatabase struct {
	data map[string]*database.UserDAO
}

func (m *mockDatabase) UpdateUser(u *database.UserDAO) error {
	m.data[u.ID] = u
	return nil
}

func (m *mockDatabase) CreateUser(u *database.UserDAO) error {
	m.data[u.ID] = u
	return nil
}

func (m *mockDatabase) GetUserByID(id string) (*database.UserDAO, error) {
	if id == "" {
		return nil, nil
	}
	v, ok := m.data[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return v, nil
}

func (m *mockDatabase) GetRoot() (*database.UserDAO, error) {
	v, ok := m.data["root"]
	if !ok {
		return nil, nil
	}
	return m.data[v.LeftID], nil
}
func (m *mockDatabase) SetRoot(id string) error {
	m.data["root"] = new(database.UserDAO)
	m.data["root"].LeftID = id
	return nil
}

func (m *mockDatabase) RemoveUserByID(id string) error {
	delete(m.data, id)
	return nil
}

var sampleTree1 = map[string]*database.UserDAO{
	"root": {
		LeftID: "1",
	},
	"1": {
		ID:         "1",
		Score:      120,
		LeftID:     "2",
		RightID:    "3",
		RightCount: 3,
	},
	"2": {
		ID:         "2",
		Score:      90,
		LeftID:     "4",
		RightID:    "5",
		RightCount: 3,
	},
	"3": {
		ID:         "3",
		Score:      140,
		LeftID:     "6",
		RightID:    "",
		RightCount: 1,
	},
	"4": {
		ID:         "4",
		Score:      70,
		LeftID:     "",
		RightID:    "",
		RightCount: 1,
	},
	"5": {
		ID:         "5",
		Score:      110,
		LeftID:     "7",
		RightID:    "",
		RightCount: 1,
	},
	"6": {
		ID:         "6",
		Score:      125,
		LeftID:     "",
		RightID:    "",
		RightCount: 1,
	},
	"7": {
		ID:         "7",
		Score:      100,
		LeftID:     "",
		RightID:    "",
		RightCount: 1,
	},
}

func TestBusiness_countNumberOfNodes(t *testing.T) {
	type fields struct {
		handler database.IDatabase
	}
	type args struct {
		root *database.UserDAO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "1",

			fields: fields{
				handler: &mockDatabase{
					data: sampleTree1,
				},
			},
			args: args{
				root: &database.UserDAO{
					ID:         "1",
					Score:      120,
					LeftID:     "2",
					RightID:    "3",
					RightCount: 3,
				},
			},
			want:    7,
			wantErr: false,
		},
		{
			name: "2",

			fields: fields{
				handler: &mockDatabase{
					data: sampleTree1,
				},
			},
			args: args{
				root: &database.UserDAO{
					ID:         "2",
					Score:      90,
					LeftID:     "4",
					RightID:    "5",
					RightCount: 3,
				},
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Business{
				handler: tt.fields.handler,
			}
			got, err := b.countNumberOfNodes(tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("Business.countNumberOfNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Business.countNumberOfNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBusiness_bstTraverse(t *testing.T) {
	type fields struct {
		handler database.IDatabase
	}
	type args struct {
		root *database.UserDAO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.UserDTOList
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{handler: &mockDatabase{
				data: sampleTree1,
			}},
			args: args{
				root: &database.UserDAO{
					ID:         "1",
					Score:      120,
					LeftID:     "2",
					RightID:    "3",
					RightCount: 3,
				},
			},
			want: models.UserDTOList{
				models.UserDTO{
					ID:    "3",
					Score: 140,
				},
				models.UserDTO{
					ID:    "6",
					Score: 125,
				},
				models.UserDTO{
					ID:    "1",
					Score: 120,
				},
				models.UserDTO{
					ID:    "5",
					Score: 110,
				},
				models.UserDTO{
					ID:    "7",
					Score: 100,
				},
				models.UserDTO{
					ID:    "2",
					Score: 90,
				},
				models.UserDTO{
					ID:    "4",
					Score: 70,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Business{
				handler: tt.fields.handler,
			}
			got, err := b.bstTraverse(tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("Business.bstTraverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Business.bstTraverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
