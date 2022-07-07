package profiles

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/internal/users"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

type ProfilesService struct {
	UsersClient users.UsersClient
	Firestore   firestore.Client
}

func NewProfilesService(usersClient users.UsersClient, firestore firestore.Client) ProfilesService {
	return ProfilesService{
		UsersClient: usersClient,
		Firestore:   firestore,
	}
}

const followsCollectionName = "follows"

type followDocData struct {
	Follower string `firestore:"follower"`
	Followee string `firestore:"followee"`
}

func newFollowDocData(follower string, followee string) followDocData {
	return followDocData{
		Follower: follower,
		Followee: followee,
	}
}

func (s *ProfilesService) Follow(ctx context.Context, followerUsername string, followeeUsername string) (*Profile, error) {
	follower, err := s.UsersClient.GetUserByUsername(followerUsername)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting follower %s", followerUsername)
		return nil, err
	}

	followee, err := s.UsersClient.GetUserByUsername(followeeUsername)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting followee %s", followeeUsername)
		return nil, err
	}

	followDocRef := s.Firestore.Collection(followsCollectionName).NewDoc()

	followdata := newFollowDocData(follower.User.Username, followee.User.Username)

	_, err = followDocRef.Create(ctx, followdata)
	if err != nil {
		return nil, err
	}

	profile, err := s.GetProfileByUsername(ctx, followeeUsername, &followerUsername)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *ProfilesService) IsFollowing(ctx context.Context, followerUsername string, followeeUsername string) (bool, error) {
	followsCollection := s.Firestore.Collection(followsCollectionName)
	query := followsCollection.Where("follower", "==", followerUsername).Where("followee", "==", followeeUsername).Limit(1)
	followsDocs := query.Documents(ctx)
	defer followsDocs.Stop()
	for {
		_, err := followsDocs.Next()
		if err == iterator.Done {
			return false, nil
		} else if err != nil {
			return false, err
		}

		return true, nil
	}
}

func (s *ProfilesService) GetProfileByUsername(ctx context.Context, username string, follower *string) (*Profile, error) {
	user, err := s.UsersClient.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	profile := NewProfile(user.User.Username, user.User.Bio, user.User.Image, nil)

	if follower != nil {
		isFollowing, err := s.IsFollowing(ctx, *follower, profile.Username)
		if err != nil {
			return nil, err
		}
		profile.Following = &isFollowing
	}

	return &profile, nil
}
