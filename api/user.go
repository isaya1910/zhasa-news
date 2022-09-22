package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/isaya1910/zhasa-news/util"
	"io/ioutil"
	"net/http"
)

type UserResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	ID        int32  `json:"id"`
	AvatarUrl string `json:"thumbnail_url"`
}

type UserRepository interface {
	GetUser(token string) (userParams CreateUserJson, err error)
}

type UserExternalRepository struct{}

func (UserExternalRepository) GetUser(token string) (userParams CreateUserJson, err error) {
	config, err := util.LoadConfig(".")
	request, err := http.NewRequest("GET", config.UserServerAddress+"/account/user/me", nil)

	if err != nil {
		return userParams, err
	}

	request.Header.Set("Authorization", token)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Print(err)
		return userParams, err
	}

	var userResponse UserResponse
	respBody, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(respBody, &userResponse)
	if err != nil {
		fmt.Print(err)
		return userParams, err
	}
	userParams.ID = &userResponse.ID
	userParams.FirstName = &userResponse.FirstName
	userParams.LastName = &userResponse.LastName
	userParams.Bio = &userResponse.Role
	userParams.AvatarUrl = &userResponse.AvatarUrl

	if *userParams.ID == 0 {
		return userParams, errors.New("user not found")
	}
	return userParams, nil
}
