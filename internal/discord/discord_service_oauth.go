package discord

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leighmacdonald/gbans/internal/domain"
	"github.com/leighmacdonald/gbans/internal/httphelper"
	"github.com/leighmacdonald/gbans/pkg/log"
)

type discordOAuthHandler struct {
	auth    domain.AuthUsecase
	config  domain.ConfigUsecase
	persons domain.PersonUsecase
	discord domain.DiscordOAuthUsecase
}

// NewHandler provides handlers for authentication with discord connect.
func NewHandler(engine *gin.Engine, auth domain.AuthUsecase, config domain.ConfigUsecase,
	persons domain.PersonUsecase, discord domain.DiscordOAuthUsecase,
) {
	handler := discordOAuthHandler{
		auth:    auth,
		config:  config,
		persons: persons,
		discord: discord,
	}

	engine.GET("/discord/oauth", handler.onOAuthDiscordCallback())

	authGrp := engine.Group("/")
	{
		// authed
		authed := authGrp.Use(auth.Middleware(domain.PUser))
		authed.GET("/api/discord/login", handler.onLogin())
		authed.GET("/api/discord/logout", handler.onLogout())
		authed.GET("/api/discord/user", handler.onGetDiscordUser())
	}
}

func (h discordOAuthHandler) onLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := httphelper.CurrentUserProfile(ctx)

		loginURL, errURL := h.discord.CreateStatefulLoginURL(currentUser.SteamID)
		if errURL != nil {
			httphelper.SetError(ctx, httphelper.NewAPIErrorf(http.StatusBadRequest, errors.Join(errURL, domain.ErrBadRequest),
				"Could not construct oauth login URL"))

			return
		}

		ctx.JSON(http.StatusOK, gin.H{"url": loginURL})
		slog.Debug("User tried to connect discord", slog.String("sid", currentUser.SteamID.String()))
	}
}

func (h discordOAuthHandler) onOAuthDiscordCallback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Query("code")
		if code == "" {
			slog.Error("Failed to get code from query")
			ctx.Redirect(http.StatusTemporaryRedirect, h.config.ExtURLRaw("/settings?section=connections"))

			return
		}

		state := ctx.Query("state")
		if state == "" {
			slog.Error("Failed to get state from query")
			ctx.Redirect(http.StatusTemporaryRedirect, h.config.ExtURLRaw("/settings?section=connections"))

			return
		}

		if err := h.discord.HandleOAuthCode(ctx, code, state); err != nil {
			slog.Error("Failed to get access token", log.ErrAttr(err))
		}

		ctx.Redirect(http.StatusTemporaryRedirect, h.config.ExtURLRaw("/settings?section=connections"))
	}
}

func (h discordOAuthHandler) onGetDiscordUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := httphelper.CurrentUserProfile(ctx)

		discord, errUser := h.discord.GetUserDetail(ctx, user.SteamID)
		if errUser != nil {
			if errors.Is(errUser, domain.ErrNoResult) {
				httphelper.SetError(ctx, httphelper.NewAPIError(http.StatusNotFound, domain.ErrNotFound))

				return
			}

			httphelper.SetError(ctx, httphelper.NewAPIErrorf(http.StatusInternalServerError, domain.ErrInternal,
				"Failed to fetch discord user details"))

			return
		}

		ctx.JSON(http.StatusOK, discord)
	}
}

func (h discordOAuthHandler) onLogout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := httphelper.CurrentUserProfile(ctx)

		errUser := h.discord.Logout(ctx, user.SteamID)
		if errUser != nil {
			if errors.Is(errUser, domain.ErrNoResult) {
				ctx.JSON(http.StatusOK, gin.H{})

				return
			}

			httphelper.SetError(ctx, httphelper.NewAPIErrorf(http.StatusInternalServerError, domain.ErrInternal,
				"Could not perform discord logout."))

			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}
