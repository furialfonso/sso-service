package team

import (
	"context"
	"cow_sso/pkg/client/restful"
	"cow_sso/pkg/config"
	"cow_sso/pkg/repository/team/dto"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	_getTeamsByUser = "/teams/user"
)

type ITeamRepository interface {
	GetTeamsByUser(ctx context.Context, userID string) (dto.TeamsByUserResponse, error)
}

type teamRepository struct {
	restfulService restful.IRestClient
}

func NewTeamRepository(restfulService restful.IRestClient) ITeamRepository {
	return &teamRepository{
		restfulService: restfulService,
	}
}

func (t *teamRepository) GetTeamsByUser(ctx context.Context, userID string) (dto.TeamsByUserResponse, error) {
	var teams dto.TeamsByUserResponse

	url := fmt.Sprintf("%s%s/%s", config.Get().UString("cow-api.url"), _getTeamsByUser, userID)
	timeOut := config.Get().UString("cow-api.timeout")
	resp, err := t.restfulService.Get(ctx, url, timeOut)
	if err != nil {
		return teams, err
	}

	err = json.Unmarshal(resp, &teams)
	if err != nil {
		return teams, errors.New("error unmarshal teams by user")
	}

	return teams, nil
}
