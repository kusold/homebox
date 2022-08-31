package v1

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/content/backend/internal/services"
	"github.com/hay-kot/content/backend/internal/types"
	"github.com/hay-kot/content/backend/pkgs/logger"
	"github.com/hay-kot/content/backend/pkgs/server"
)

// HandleUserSelf godoc
// @Summary  Get the current user
// @Tags     User
// @Produce  json
// @Param    payload  body      types.UserRegistration  true  "User Data"
// @Success  200      {object}  server.Result{item=ent.User}
// @Router   /v1/users/register [Post]
func (ctrl *V1Controller) HandleUserRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regData := types.UserRegistration{}

		if err := server.Decode(r, &regData); err != nil {
			ctrl.log.Error(err, nil)
			server.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		_, err := ctrl.svc.User.RegisterUser(r.Context(), regData)
		if err != nil {
			ctrl.log.Error(err, nil)
			server.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		server.Respond(w, http.StatusNoContent, nil)
	}
}

// HandleUserSelf godoc
// @Summary   Get the current user
// @Tags      User
// @Produce   json
// @Success   200  {object}  server.Result{item=ent.User}
// @Router    /v1/users/self [GET]
// @Security  Bearer
func (ctrl *V1Controller) HandleUserSelf() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := services.UseTokenCtx(r.Context())
		usr, err := ctrl.svc.User.GetSelf(r.Context(), token)
		if usr.ID == uuid.Nil || err != nil {
			ctrl.log.Error(errors.New("no user within request context"), nil)
			server.RespondServerError(w)
			return
		}

		server.Respond(w, http.StatusOK, server.Wrap(usr))
	}
}

// HandleUserUpdate godoc
// @Summary   Update the current user
// @Tags      User
// @Produce   json
// @Param     payload  body      types.UserUpdate  true  "User Data"
// @Success   200      {object}  server.Result{item=types.UserUpdate}
// @Router    /v1/users/self [PUT]
// @Security  Bearer
func (ctrl *V1Controller) HandleUserUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updateData := types.UserUpdate{}
		if err := server.Decode(r, &updateData); err != nil {
			ctrl.log.Error(err, logger.Props{
				"scope":   "user",
				"details": "failed to decode user update data",
			})
			server.RespondError(w, http.StatusBadRequest, err)
			return
		}

		actor := services.UseUserCtx(r.Context())
		newData, err := ctrl.svc.User.UpdateSelf(r.Context(), actor.ID, updateData)

		if err != nil {
			ctrl.log.Error(err, logger.Props{
				"scope":   "user",
				"details": "failed to update user",
			})
			server.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		server.Respond(w, http.StatusOK, server.Wrap(newData))
	}
}

// HandleUserUpdatePassword godoc
// @Summary   Update the current user's password // TODO:
// @Tags      User
// @Produce   json
// @Success   204
// @Router    /v1/users/self/password [PUT]
// @Security  Bearer
func (ctrl *V1Controller) HandleUserUpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
