package profiles

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Profile struct {
	Profile profileProfile `json:"profile"`
}

type profileProfile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func FollowUser(username string, token string) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:8080/profiles/%s/follow", username)

	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func FollowUserAndDecode(username string, token string) (*Profile, error) {
	response, err := FollowUser(username, token)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("got %d, want %d", response.StatusCode, http.StatusCreated)
	}

	responseData := &Profile{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
