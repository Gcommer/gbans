package srcds

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leighmacdonald/gbans/internal/domain"
	"github.com/leighmacdonald/gbans/internal/httphelper"
)

type speedrunHandler struct {
	speedruns domain.SpeedrunUsecase
	auth      domain.AuthUsecase
	config    domain.ConfigUsecase
}

func NewHandler(engine *gin.Engine, speedruns domain.SpeedrunUsecase, auth domain.AuthUsecase, config domain.ConfigUsecase) {
	handler := speedrunHandler{
		speedruns: speedruns,
		auth:      auth,
		config:    config,
	}

	guestGroup := engine.Group("/")
	{
		guest := guestGroup.Use(auth.Middleware(domain.PGuest))
		// Groups
		// guest.GET("/api/speedruns/overall", handler.getOverall())
		guest.GET("/api/speedruns/map", handler.getByMap())
		guest.GET("/api/speedruns/overall/top", handler.getOverallTopN())
		guest.GET("/api/speedruns/overall/recent", handler.getRecentChanges())
		guest.GET("/api/speedruns/byid/:speedrun_id", handler.getSpeedrun())
	}

	srcdsGroup := engine.Group("/")
	{
		server := srcdsGroup.Use(auth.MiddlewareServer())
		server.POST("/api/sm/speedruns", handler.postSpeedrun())
	}
}

func (s *speedrunHandler) postSpeedrun() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var sr domain.Speedrun
		if !httphelper.Bind(ctx, &sr) {
			return
		}

		speedrun, errSpeedrun := s.speedruns.Save(ctx, sr)
		if errSpeedrun != nil {
			httphelper.SetError(ctx, httphelper.NewAPIError(http.StatusInternalServerError, errors.Join(errSpeedrun, domain.ErrInternal)))

			return
		}

		ctx.JSON(http.StatusOK, speedrun)
	}
}

// func (s *speedrunHandler) getOverall() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		top, errTop := s.speedruns.TopNOverall(ctx, 3)
//		if errTop != nil {
//			slog.Error("Failed to load top speedruns", errTop)
//			httphelper.HandleErrInternal(ctx)
//
//			return
//		}
//
//		ctx.JSON(http.StatusOK, top)
//	}
// }

func (s *speedrunHandler) getByMap() gin.HandlerFunc {
	type q struct {
		MapName string `schema:"map_name"`
	}

	return func(ctx *gin.Context) {
		var query q
		if !httphelper.BindQuery(ctx, &query) {
			return
		}

		runs, errRuns := s.speedruns.ByMap(ctx, query.MapName)
		if errRuns != nil {
			httphelper.SetError(ctx, httphelper.NewAPIError(http.StatusInternalServerError, errors.Join(errRuns, domain.ErrInternal)))

			return
		}

		ctx.JSON(http.StatusOK, runs)
	}
}

func (s *speedrunHandler) getSpeedrun() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		speedrunID, idFound := httphelper.GetIntParam(ctx, "speedrun_id")
		if !idFound {
			return
		}

		speedrun, errSpeedrun := s.speedruns.ByID(ctx, speedrunID)
		if errSpeedrun != nil {
			if errors.Is(errSpeedrun, domain.ErrNoResult) {
				httphelper.SetError(ctx, httphelper.NewAPIError(http.StatusNotFound, domain.ErrNotFound))

				return
			}

			httphelper.SetError(ctx, httphelper.NewAPIError(http.StatusInternalServerError, errors.Join(errSpeedrun, domain.ErrInternal)))

			return
		}

		ctx.JSON(http.StatusOK, speedrun)
	}
}

func (s *speedrunHandler) getRecentChanges() gin.HandlerFunc {
	var query struct {
		Count int `json:"count"`
	}

	return func(ctx *gin.Context) {
		if !httphelper.BindQuery(ctx, &query) {
			return
		}

		top, errTop := s.speedruns.Recent(ctx, query.Count)
		if errTop != nil {
			if errors.Is(errTop, domain.ErrValueOutOfRange) {
				httphelper.SetError(ctx, httphelper.NewAPIError(http.StatusBadRequest, domain.ErrValueOutOfRange))

				return
			}

			httphelper.SetError(ctx, httphelper.NewAPIError(http.StatusInternalServerError, errors.Join(errTop, domain.ErrInternal)))

			return
		}

		ctx.JSON(http.StatusOK, top)
	}
}

func (s *speedrunHandler) getOverallTopN() gin.HandlerFunc {
	var query struct {
		Count int `json:"count"`
	}

	return func(ctx *gin.Context) {
		if !httphelper.BindQuery(ctx, &query) {
			return
		}

		top, errTop := s.speedruns.TopNOverall(ctx, query.Count)
		if errTop != nil {
			if errors.Is(errTop, domain.ErrValueOutOfRange) {
				httphelper.SetError(ctx, httphelper.NewAPIError(http.StatusBadRequest, domain.ErrValueOutOfRange))

				return
			}

			httphelper.SetError(ctx, httphelper.NewAPIError(http.StatusInternalServerError, errors.Join(errTop, domain.ErrInternal)))

			return
		}

		ctx.JSON(http.StatusOK, top)
	}
}
