package profiles

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/internal/auth"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/internal/custom_errors"
	"github.com/rs/zerolog/log"
)

type ProfilesHandlers struct {
	ProfilesService ProfilesService
}

func NewProfilesHandlers(profilesService ProfilesService) ProfilesHandlers {
	return ProfilesHandlers{
		ProfilesService: profilesService,
	}
}

type profileResponse struct {
	Profile profileResponseProfile `json:"profile"`
}

type profileResponseProfile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func newProfileResponse(username string, bio string, image string, following bool) profileResponse {
	return profileResponse{
		Profile: profileResponseProfile{
			Username:  username,
			Bio:       bio,
			Image:     image,
			Following: following,
		},
	}
}

type errorResponse struct {
	Errors errorResponseErrors `json:"errors"`
}

type errorResponseErrors struct {
	Body []string `json:"body"`
}

func newErrorResponse(errors []error) errorResponse {
	var body []string
	for _, err := range errors {
		body = append(body, err.Error())
	}

	return errorResponse{
		Errors: errorResponseErrors{
			Body: body,
		},
	}
}

func (h *ProfilesHandlers) FollowUser(w http.ResponseWriter, r *http.Request) {
	follower := r.Context().Value(auth.UsernameContextKey).(string)
	followee := chi.URLParam(r, "username")

	profile, err := h.ProfilesService.Follow(r.Context(), follower, followee)
	if err != nil {
		if _, ok := err.(*custom_errors.NotFoundError); ok {
			notFound(w, r, []error{err})
			return
		}

		if _, ok := err.(*custom_errors.AlreadyExistsError); ok {
			unprocessableEntity(w, r, []error{err})
			return
		}

		internalServerError(w, r, err)
		return
	}

	responseBody := newProfileResponse(profile.Username, profile.Bio, profile.Image, *profile.Following)

	response, err := json.Marshal(responseBody)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *ProfilesHandlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	follower := r.Context().Value(auth.UsernameContextKey).(string)
	followee := chi.URLParam(r, "username")

	profile, err := h.ProfilesService.GetProfileByUsername(r.Context(), followee, &follower)
	if err != nil {
		if _, ok := err.(*custom_errors.NotFoundError); ok {
			notFound(w, r, []error{err})
			return
		}

		if _, ok := err.(*custom_errors.AlreadyExistsError); ok {
			unprocessableEntity(w, r, []error{err})
			return
		}

		internalServerError(w, r, err)
		return
	}

	responseBody := newProfileResponse(profile.Username, profile.Bio, profile.Image, *profile.Following)

	response, err := json.Marshal(responseBody)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func notFound(w http.ResponseWriter, r *http.Request, errors []error) {
	response, err := json.Marshal(newErrorResponse(errors))
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(response)
}

func unprocessableEntity(w http.ResponseWriter, r *http.Request, errors []error) {
	response, err := json.Marshal(newErrorResponse(errors))
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(response)
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Error().Err(err).Msg("")
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
