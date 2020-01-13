package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"image"
	"net/http"

	"github.com/git-sim/tc/app/domain/repo"

	"github.com/git-sim/tc/app/usecase"
)

// ProfileCtxFunc returns a context handler to validate the accountID exists, and is registered
func ProfileCtxFunc(pu usecase.ProfileUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accountID, ok := getAccountIDFromContext(w, r)
			if !ok {
				return
			}

			profile, err := pu.Get(uint64(accountID))
			if err != nil || profile == nil {
				ReportUsecaseFault(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), "profile", profile)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetProfile ...
func GetProfile(pu usecase.ProfileUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID, ok := getAccountIDFromContext(w, r)
		if !ok {
			return
		}

		prof, err := pu.Get(uint64(accountID))
		if err != nil {
			ReportUsecaseFault(w, err)
			return
		}

		err = json.NewEncoder(w).Encode(prof)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
	}
}

// PutProfile ...
func PutProfile(pu usecase.ProfileUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The message is already cached in the context for this request
		// extract, modify, update  (concurrency?)
		accountID, ok := getAccountIDFromContext(w, r)
		if !ok {
			return
		}

		// Decode body
		var inprofile repo.PublicProfile
		err := decodeJSONBody(w, r, &inprofile)
		if err != nil {
			var es *ErrorJSONDecode
			if errors.As(err, &es) {
				http.Error(w, es.msg, es.status)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}

		err = pu.Set(uint64(accountID), &inprofile)
		if err != nil {
			ReportUsecaseFault(w, err)
			return
		}

		outprofile, err := pu.Get(uint64(accountID))
		if err != nil {
			ReportUsecaseFault(w, err)
			return
		}

		err = json.NewEncoder(w).Encode(outprofile)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
	}
}

type profileJustBio struct {
	Bio string `json:"bio"`
}

// GetProfileBio ...
func GetProfileBio(pu usecase.ProfileUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, ok := getProfileFromContext(w, r)
		if !ok {
			return
		}
		out := profileJustBio{Bio: profile.NameAndBios[repo.EnumBio]}
		err := json.NewEncoder(w).Encode(&out)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
	}
}

// PutProfileBio ...
func PutProfileBio(pu usecase.ProfileUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID, ok := getAccountIDFromContext(w, r)
		if !ok {
			return
		}
		profile, ok := getProfileFromContext(w, r)
		if !ok {
			return
		}

		// Decode body
		var inprofile profileJustBio
		err := decodeJSONBody(w, r, &inprofile)
		if err != nil {
			var es *ErrorJSONDecode
			if errors.As(err, &es) {
				http.Error(w, es.msg, es.status)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		profile.NameAndBios[repo.EnumBio] = inprofile.Bio
		pu.Set(uint64(accountID), profile)
	}
}

type profileJustAvatar struct {
	Avatar image.Image `json:"avatar"`
}

// GetProfileAvatar ...
func GetProfileAvatar(pu usecase.ProfileUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, ok := getProfileFromContext(w, r)
		if !ok {
			return
		}

		out := profileJustAvatar{}
		if profile.Pics[repo.EnumAvatar] != nil {
			out.Avatar = *profile.Pics[repo.EnumAvatar]
		}

		err := json.NewEncoder(w).Encode(&out)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
	}
}

// PutProfileAvatar ...
func PutProfileAvatar(pu usecase.ProfileUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID, ok := getAccountIDFromContext(w, r)
		if !ok {
			return
		}
		profile, ok := getProfileFromContext(w, r)
		if !ok {
			return
		}

		// Decode body
		var inprofile profileJustAvatar
		err := decodeJSONBody(w, r, &inprofile)
		if err != nil {
			var es *ErrorJSONDecode
			if errors.As(err, &es) {
				http.Error(w, es.msg, es.status)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		profile.Pics[repo.EnumAvatar] = &inprofile.Avatar
		pu.Set(uint64(accountID), profile)
	}
}

type profileJustBg struct {
	Background image.Image `json:"background"`
}

// GetProfileBackground ...
func GetProfileBackground(pu usecase.ProfileUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, ok := getProfileFromContext(w, r)
		if !ok {
			return
		}
		out := profileJustBg{}
		if profile.Pics[repo.EnumBackground] != nil {
			out.Background = *profile.Pics[repo.EnumBackground]
		}

		err := json.NewEncoder(w).Encode(&out)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
	}
}

// PutProfileBackground ...
func PutProfileBackground(pu usecase.ProfileUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID, ok := getAccountIDFromContext(w, r)
		if !ok {
			return
		}
		profile, ok := getProfileFromContext(w, r)
		if !ok {
			return
		}

		// Decode body
		var inprofile profileJustBg
		err := decodeJSONBody(w, r, &inprofile)
		if err != nil {
			var es *ErrorJSONDecode
			if errors.As(err, &es) {
				http.Error(w, es.msg, es.status)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		profile.Pics[repo.EnumBackground] = &inprofile.Background
		pu.Set(uint64(accountID), profile)
	}
}
