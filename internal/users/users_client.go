package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/internal/custom_errors"
)

type UsersClient struct {
	BaseURL string
}

func NewUsersClient(baseURL string) UsersClient {
	return UsersClient{
		BaseURL: baseURL,
	}
}

type GetCurrentUserResponse struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Token    string `json:"token"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

type GetUserByUsernameResponse struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

type errorResponse struct {
	Errors *errorResponseErrors `json:"errors"`
}

type errorResponseErrors struct {
	Body []string `json:"body"`
}

func (c *UsersClient) GetCurrentUser(tokenString string) (*GetCurrentUserResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/user", c.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusUnprocessableEntity:
		responseData := &errorResponse{}
		err = json.NewDecoder(response.Body).Decode(&responseData)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(responseData.Errors.Body[0])
	default:
		return nil, errors.New(fmt.Sprintf("GetCurrentUser error. HTTP status %d", response.StatusCode))
	}

	responseData := GetCurrentUserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

func (c *UsersClient) GetUserByUsername(username string) (*GetUserByUsernameResponse, error) {
	url := fmt.Sprintf("%s/users/%s", c.BaseURL, username)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, &custom_errors.NotFoundError{
			Message: "User not found",
		}
	case http.StatusUnprocessableEntity:
		responseData := &errorResponse{}
		err = json.NewDecoder(response.Body).Decode(&responseData)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(responseData.Errors.Body[0])
	default:
		return nil, errors.New(fmt.Sprintf("GetUserByUsername error. HTTP status %d", response.StatusCode))
	}

	responseData := GetUserByUsernameResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}
