package profiles

import (
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-client/pkg/client"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/test/utils"
)

func TestValidRequestWhenFollowUserShouldReturnProfile(t *testing.T) {
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
