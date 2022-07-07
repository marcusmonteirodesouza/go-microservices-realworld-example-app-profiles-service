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

func TestGivenUserIsFollowedWhenGetProfileShouldReturnProfileWithFollowingTrue(t *testing.T) {
	profileClient := utils.NewApiClient()

	profileUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	profileEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	profilePassword := faker.Password()
	profileBio := faker.Paragraph()
	profileImage := faker.URL()

	profileUser, err := profileClient.Users.RegisterUser(profileUsername, profileEmail, profilePassword)
	if err != nil {
		t.Fatal(err)
	}

	profileUser, err = profileClient.Users.UpdateUser(client.UpdateUserRequest{
		Bio:   &profileBio,
		Image: &profileImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	callerClient := utils.NewApiClient()

	callerUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	callerEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	callerPassword := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Password())

	caller, err := callerClient.Users.RegisterUser(callerUsername, callerEmail, callerPassword)
	if err != nil {
		t.Fatal(err)
	}

	profile, err := FollowUserAndDecode(profileUser.User.Username, caller.User.Token)
	if err != nil {
		t.Fatal(err)
	}

	profile, err = GetProfileAndDecode(profileUser.User.Username, &caller.User.Token)
	if err != nil {
		t.Fatal(err)
	}

	if profile.Profile.Username != profileUser.User.Username {
		t.Fatalf("got %s, want %s", profile.Profile.Username, profileUser.User.Username)
	}

	if profile.Profile.Bio != profileUser.User.Bio {
		t.Fatalf("got %s, want %s", profile.Profile.Bio, profileUser.User.Bio)
	}

	if profile.Profile.Image != profileUser.User.Image {
		t.Fatalf("got %s, want %s", profile.Profile.Image, profileUser.User.Image)
	}

	if !profile.Profile.Following {
		t.Fatalf("got %t, want %t", profile.Profile.Following, true)
	}
}

func TestGivenUserIsNotFollowedWhenGetProfileShouldReturnProfileWithFollowingFalse(t *testing.T) {
	profileClient := utils.NewApiClient()

	profileUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	profileEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	profilePassword := faker.Password()
	profileBio := faker.Paragraph()
	profileImage := faker.URL()

	profileUser, err := profileClient.Users.RegisterUser(profileUsername, profileEmail, profilePassword)
	if err != nil {
		t.Fatal(err)
	}

	profileUser, err = profileClient.Users.UpdateUser(client.UpdateUserRequest{
		Bio:   &profileBio,
		Image: &profileImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	callerClient := utils.NewApiClient()

	callerUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	callerEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	callerPassword := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Password())

	caller, err := callerClient.Users.RegisterUser(callerUsername, callerEmail, callerPassword)
	if err != nil {
		t.Fatal(err)
	}

	profile, err := GetProfileAndDecode(profileUser.User.Username, &caller.User.Token)
	if err != nil {
		t.Fatal(err)
	}

	if profile.Profile.Username != profileUser.User.Username {
		t.Fatalf("got %s, want %s", profile.Profile.Username, profileUser.User.Username)
	}

	if profile.Profile.Bio != profileUser.User.Bio {
		t.Fatalf("got %s, want %s", profile.Profile.Bio, profileUser.User.Bio)
	}

	if profile.Profile.Image != profileUser.User.Image {
		t.Fatalf("got %s, want %s", profile.Profile.Image, profileUser.User.Image)
	}

	if profile.Profile.Following {
		t.Fatalf("got %t, want %t", profile.Profile.Following, false)
	}
}

func TestGivenInvalidTokenWhenGetProfileShouldReturnProfileWithFollowingFalse(t *testing.T) {
	profileClient := utils.NewApiClient()

	profileUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	profileEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	profilePassword := faker.Password()
	profileBio := faker.Paragraph()
	profileImage := faker.URL()

	profile, err := profileClient.Users.RegisterUser(profileUsername, profileEmail, profilePassword)
	if err != nil {
		t.Fatal(err)
	}

	profile, err = profileClient.Users.UpdateUser(client.UpdateUserRequest{
		Bio:   &profileBio,
		Image: &profileImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	callerClient := utils.NewApiClient()

	callerUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	callerEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	callerPassword := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Password())

	_, err = callerClient.Users.RegisterUser(callerUsername, callerEmail, callerPassword)
	if err != nil {
		t.Fatal(err)
	}

	invalidToken := "invalidToken"

	response, err := GetProfile(profile.User.Username, &invalidToken)
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

func TestGivenUserIsNotFoundWhenGetProfileShouldReturnNotFound(t *testing.T) {
	profileUsername := faker.UUIDHyphenated()

	callerClient := utils.NewApiClient()

	callerUsername := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Username())
	callerEmail := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Email())
	callerPassword := fmt.Sprintf("%s%s", utils.TestPrefix, faker.Password())

	caller, err := callerClient.Users.RegisterUser(callerUsername, callerEmail, callerPassword)
	if err != nil {
		t.Fatal(err)
	}

	response, err := FollowUser(profileUsername, caller.User.Token)
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
