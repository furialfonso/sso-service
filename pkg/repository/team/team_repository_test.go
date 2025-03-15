package team

import (
	"context"
	"cow_sso/mocks"
	"cow_sso/pkg/repository/team/dto"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTeamRepository struct {
	restfulService *mocks.IRestClient
}

type teamMocks struct {
	teamRepository func(f *mockTeamRepository)
}

func Test_GetTeamsByUser(t *testing.T) {
	tests := []struct {
		name   string
		userID string
		mocks  teamMocks
		outPut dto.TeamsByUserResponse
		expErr error
	}{
		{
			name:   "error client",
			userID: "ABC",
			mocks: teamMocks{
				func(f *mockTeamRepository) {
					f.restfulService.On("Get", mock.Anything, mock.Anything, "5s").Return(nil, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:   "error unmarshal teams by user",
			userID: "ABC",
			mocks: teamMocks{
				func(f *mockTeamRepository) {
					b, _ := json.Marshal("teams")
					f.restfulService.On("Get", mock.Anything, mock.Anything, "5s").Return(b, nil)
				},
			},
			expErr: errors.New("error unmarshal teams by user"),
		},
		{
			name:   "full flow",
			userID: "ABC",
			mocks: teamMocks{
				func(f *mockTeamRepository) {
					teams := dto.TeamsByUserResponse{
						Teams: []dto.TeamResponse{
							{
								Code: "XYZ",
								Debt: 100,
							},
						},
					}
					b, _ := json.Marshal(teams)
					f.restfulService.On("Get", mock.Anything, mock.Anything, "5s").Return(b, nil)
				},
			},
			outPut: dto.TeamsByUserResponse{
				Teams: []dto.TeamResponse{
					{
						Code: "XYZ",
						Debt: 100,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockTeamRepository{
				restfulService: &mocks.IRestClient{},
			}
			tt.mocks.teamRepository(m)
			r := NewTeamRepository(m.restfulService)
			team, err := r.GetTeamsByUser(context.Background(), tt.userID)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
			assert.Equal(t, tt.outPut, team)
		})
	}
}
