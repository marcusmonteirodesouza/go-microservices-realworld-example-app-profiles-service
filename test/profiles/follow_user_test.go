package profiles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-client/pkg/client"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/test/utils"
)

func TestGivenValidRequestWhenFollowUserShouldReturnProfile(t *testing.T) {
	followeeClient := utils.NewApiClient()

	followeeUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	followeeEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	followeePassword := faker.Password()
	followeeBio := faker.Paragraph()
	followeeImage := faker.URL()

	followee, err := followeeClient.Users.RegisterUser(followeeUsername, followeeEmail, followeePassword)
	if err != nil {
		t.Fatal(err)
	}

	followee, err = followeeClient.Users.UpdateUser(client.UpdateUserRequest{
		Bio:   &followeeBio,
		Image: &followeeImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	followerClient := utils.NewApiClient()

	followerUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	followerEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	followerPassword := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Password())

	follower, err := followerClient.Users.RegisterUser(followerUsername, followerEmail, followerPassword)
	if err != nil {
		t.Fatal(err)
	}

	profile, err := FollowUserAndDecode(followee.User.Username, follower.User.Token)
	if err != nil {
		t.Fatal(err)
	}

	if profile.Profile.Username != followee.User.Username {
		t.Fatalf("got %s, want %s", profile.Profile.Username, followee.User.Username)
	}

	if profile.Profile.Bio != followee.User.Bio {
		t.Fatalf("got %s, want %s", profile.Profile.Bio, followee.User.Bio)
	}

	if profile.Profile.Image != followee.User.Image {
		t.Fatalf("got %s, want %s", profile.Profile.Image, followee.User.Image)
	}

	if !profile.Profile.Following {
		t.Fatalf("got %t, want %t", profile.Profile.Following, true)
	}
}

func TestGivenInvalidAccessTokenWhenFollowUserShouldReturnUnauthorized(t *testing.T) {
	followeeClient := utils.NewApiClient()

	followeeUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	followeeEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	followeePassword := faker.Password()
	followeeBio := faker.Paragraph()
	followeeImage := faker.URL()

	followee, err := followeeClient.Users.RegisterUser(followeeUsername, followeeEmail, followeePassword)
	if err != nil {
		t.Fatal(err)
	}

	followee, err = followeeClient.Users.UpdateUser(client.UpdateUserRequest{
		Bio:   &followeeBio,
		Image: &followeeImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	followerClient := utils.NewApiClient()

	followerUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	followerEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	followerPassword := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Password())

	_, err = followerClient.Users.RegisterUser(followerUsername, followerEmail, followerPassword)
	if err != nil {
		t.Fatal(err)
	}

	response, err := FollowUser(followee.User.Username, "invalidToken")
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	bodyString := strings.TrimSpace(string(bodyBytes))
	if bodyString != "Unauthorized" {
		t.Fatalf("got %s, want %s", bodyString, "Unauthorized")
	}
}

func TestGivenFolloweeIsNotFoundWhenFollowUserShouldReturnNotFound(t *testing.T) {
	followeeUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())

	followerClient := utils.NewApiClient()

	followerUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	followerEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	followerPassword := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Password())

	follower, err := followerClient.Users.RegisterUser(followerUsername, followerEmail, followerPassword)
	if err != nil {
		t.Fatal(err)
	}

	response, err := FollowUser(followeeUsername, follower.User.Token)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusNotFound)
	}

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "User not found" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "User not found")
	}
}
