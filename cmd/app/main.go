package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/internal/auth"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/internal/firestore"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/internal/profiles"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/internal/users"
	"github.com/rs/zerolog/log"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal().Err(err).Msg("Environment variable 'PORT' must be set and set to an integer")
	}

	firestoreProjectId := os.Getenv("FIRESTORE_PROJECT_ID")
	if len(firestoreProjectId) == 0 {
		log.Fatal().Err(err).Msg("Environment variable 'FIRESTORE_PROJECT_ID' must be set and not be empty")
	}

	ctx := context.Background()

	usersServiceBaseUrl := os.Getenv("API_BASE_URL")
	if len(usersServiceBaseUrl) == 0 {
		log.Fatal().Err(err).Msg("Environment variable 'API_BASE_URL' must be set and not be empty")
	}

	firestoreClient, err := firestore.InitFirestore(ctx, firestoreProjectId)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing the Firestore client")
	}

	defer firestoreClient.Close()

	usersClient := users.NewUsersClient(usersServiceBaseUrl)

	profilesService := profiles.NewProfilesService(usersClient, *firestoreClient)

	profilesHandlers := profiles.NewProfilesHandlers(profilesService)

	authMiddleware := auth.NewAuthMiddleware(usersClient)

	router := chi.NewRouter()
	router.Post("/profiles/{username}/follow", authMiddleware.Authenticate(profilesHandlers.FollowUser))
	router.Get("/profiles/{username}", authMiddleware.AuthenticateOptional(profilesHandlers.GetProfile))
	router.Delete("/profiles/{username}/follow", authMiddleware.Authenticate(profilesHandlers.UnfollowUser))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	log.Info().Msgf("Starting server on %s", server.Addr)
	log.Fatal().Err(server.ListenAndServe()).Msg("")
}
