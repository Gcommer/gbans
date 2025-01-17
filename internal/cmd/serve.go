package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leighmacdonald/gbans/internal/app"
	"github.com/leighmacdonald/gbans/internal/appeal"
	"github.com/leighmacdonald/gbans/internal/asset"
	"github.com/leighmacdonald/gbans/internal/auth"
	"github.com/leighmacdonald/gbans/internal/ban"
	"github.com/leighmacdonald/gbans/internal/blocklist"
	"github.com/leighmacdonald/gbans/internal/chat"
	"github.com/leighmacdonald/gbans/internal/config"
	"github.com/leighmacdonald/gbans/internal/contest"
	"github.com/leighmacdonald/gbans/internal/database"
	"github.com/leighmacdonald/gbans/internal/demo"
	"github.com/leighmacdonald/gbans/internal/discord"
	"github.com/leighmacdonald/gbans/internal/domain"
	"github.com/leighmacdonald/gbans/internal/forum"
	"github.com/leighmacdonald/gbans/internal/httphelper"
	"github.com/leighmacdonald/gbans/internal/match"
	"github.com/leighmacdonald/gbans/internal/metrics"
	"github.com/leighmacdonald/gbans/internal/network"
	"github.com/leighmacdonald/gbans/internal/news"
	"github.com/leighmacdonald/gbans/internal/notification"
	"github.com/leighmacdonald/gbans/internal/patreon"
	"github.com/leighmacdonald/gbans/internal/person"
	"github.com/leighmacdonald/gbans/internal/queue"
	"github.com/leighmacdonald/gbans/internal/report"
	"github.com/leighmacdonald/gbans/internal/servers"
	"github.com/leighmacdonald/gbans/internal/srcds"
	"github.com/leighmacdonald/gbans/internal/state"
	"github.com/leighmacdonald/gbans/internal/steamgroup"
	"github.com/leighmacdonald/gbans/internal/votes"
	"github.com/leighmacdonald/gbans/internal/wiki"
	"github.com/leighmacdonald/gbans/internal/wordfilter"
	"github.com/leighmacdonald/gbans/pkg/fp"
	"github.com/leighmacdonald/gbans/pkg/log"
	"github.com/leighmacdonald/gbans/pkg/logparse"
	"github.com/leighmacdonald/steamid/v4/steamid"
	"github.com/riverqueue/river"
	"github.com/spf13/cobra"
)

func firstTimeSetup(ctx context.Context, persons domain.PersonUsecase, news domain.NewsUsecase,
	wiki domain.WikiUsecase, conf domain.Config,
) error {
	_, errRootUser := persons.GetPersonBySteamID(ctx, steamid.New(conf.Owner))
	if errRootUser == nil {
		return nil
	}

	if !errors.Is(errRootUser, domain.ErrNoResult) {
		return errRootUser
	}

	newOwner := domain.NewPerson(steamid.New(conf.Owner))
	newOwner.PermissionLevel = domain.PAdmin

	if errSave := persons.SavePerson(ctx, &newOwner); errSave != nil {
		slog.Error("Failed create new owner", log.ErrAttr(errSave))
	}

	newsEntry := domain.NewsEntry{
		Title:       "Welcome to gbans",
		BodyMD:      "This is an *example* **news** entry.",
		IsPublished: true,
		CreatedOn:   time.Now(),
		UpdatedOn:   time.Now(),
	}

	if errSave := news.SaveNewsArticle(ctx, &newsEntry); errSave != nil {
		return errSave
	}

	page := domain.WikiPage{
		Slug:      domain.RootSlug,
		BodyMD:    "# Welcome to the wiki",
		Revision:  1,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}

	_, errSave := wiki.SaveWikiPage(ctx, newOwner, page.Slug, page.BodyMD, page.PermissionLevel)
	if errSave != nil {
		slog.Error("Failed save example wiki entry", log.ErrAttr(errSave))
	}

	return nil
}

func createQueueWorkers(people domain.PersonUsecase, notifications domain.NotificationUsecase,
	discordUC domain.DiscordUsecase, authRepo domain.AuthRepository,
	patreonUC domain.PatreonUsecase, reports domain.ReportUsecase, discordOAuth domain.DiscordOAuthUsecase,
) *river.Workers {
	workers := river.NewWorkers()

	river.AddWorker[notification.SenderArgs](workers, notification.NewSenderWorker(people, notifications, discordUC))
	river.AddWorker[auth.CleanupArgs](workers, auth.NewCleanupWorker(authRepo))
	river.AddWorker[patreon.AuthUpdateArgs](workers, patreon.NewSyncWorker(patreonUC))
	river.AddWorker[report.MetaInfoArgs](workers, report.NewMetaInfoWorker(reports))
	river.AddWorker[discord.TokenRefreshArgs](workers, discord.NewTokenRefreshWorker(discordOAuth))

	return workers
}

func createPeriodicJobs() []*river.PeriodicJob {
	jobs := []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(24*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return auth.CleanupArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true}),

		river.NewPeriodicJob(
			river.PeriodicInterval(time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return patreon.AuthUpdateArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true}),

		river.NewPeriodicJob(
			river.PeriodicInterval(24*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return report.MetaInfoArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true}),

		river.NewPeriodicJob(
			river.PeriodicInterval(time.Hour*12),
			func() (river.JobArgs, *river.InsertOpts) {
				return discord.TokenRefreshArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true}),
	}

	return jobs
}

// serveCmd represents the serve command.
func serveCmd() *cobra.Command { //nolint:maintidx
	return &cobra.Command{
		Use:   "serve",
		Short: "Starts the gbans web app",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			slog.Info("Starting gbans...",
				slog.String("version", app.BuildVersion),
				slog.String("commit", app.BuildCommit),
				slog.String("date", app.BuildDate))

			staticConfig, errStatic := config.ReadStaticConfig()
			if errStatic != nil {
				slog.Error("Failed to read static config", log.ErrAttr(errStatic))

				return errStatic
			}

			dbConn := database.New(staticConfig.DatabaseDSN, staticConfig.DatabaseAutoMigrate, staticConfig.DatabaseLogQueries)
			if errConnect := dbConn.Connect(ctx); errConnect != nil {
				slog.Error("Cannot initialize database", log.ErrAttr(errConnect))

				return errConnect
			}

			defer func() {
				if errClose := dbConn.Close(); errClose != nil {
					slog.Error("Failed to close database cleanly")
				}
			}()

			if err := queue.Init(ctx, dbConn.Pool()); err != nil {
				slog.Error("Failed to initialize queue", log.ErrAttr(err))

				return err
			}

			// Config
			configUsecase := config.NewConfigUsecase(staticConfig, config.NewConfigRepository(dbConn))
			if err := configUsecase.Init(ctx); err != nil {
				slog.Error("Failed to init config", log.ErrAttr(err))

				return err
			}

			if errConfig := configUsecase.Reload(ctx); errConfig != nil {
				slog.Error("Failed to read config", log.ErrAttr(errConfig))

				return errConfig
			}

			conf := configUsecase.Config()

			if conf.Sentry.SentryDSN != "" {
				sentryClient, err := log.NewSentryClient(conf.Sentry.SentryDSN, conf.Sentry.SentryTrace, conf.Sentry.SentrySampleRate, app.BuildVersion)
				if err != nil {
					slog.Error("Failed to setup sentry client")
				} else {
					defer sentryClient.Flush(2 * time.Second)
				}
			}

			logCloser := log.MustCreateLogger(conf.Log.File, conf.Log.Level)
			defer logCloser()

			eventBroadcaster := fp.NewBroadcaster[logparse.EventType, logparse.ServerEvent]()
			weaponsMap := fp.NewMutexMap[logparse.Weapon, int]()

			discordRepository, errDR := discord.NewDiscordRepository(conf)
			if errDR != nil {
				slog.Error("Cannot initialize discord", log.ErrAttr(errDR))

				return errDR
			}

			discordUsecase := discord.NewDiscordUsecase(discordRepository, configUsecase)

			if err := discordUsecase.Start(); err != nil {
				slog.Error("Failed to start discord", log.ErrAttr(err))

				return err
			}

			notificationUsecase := notification.NewNotificationUsecase(notification.NewNotificationRepository(dbConn), discordUsecase)

			wordFilterUsecase := wordfilter.NewWordFilterUsecase(wordfilter.NewWordFilterRepository(dbConn), notificationUsecase)
			if err := wordFilterUsecase.Import(ctx); err != nil {
				slog.Error("Failed to load word filters", log.ErrAttr(err))

				return err
			}

			defer discordUsecase.Shutdown(conf.Discord.GuildID)

			personUsecase := person.NewPersonUsecase(person.NewPersonRepository(conf, dbConn), configUsecase)

			networkUsecase := network.NewNetworkUsecase(eventBroadcaster, network.NewNetworkRepository(dbConn), personUsecase, configUsecase)
			go networkUsecase.Start(ctx)

			assetRepository := asset.NewLocalRepository(dbConn, configUsecase)
			if err := assetRepository.Init(ctx); err != nil {
				slog.Error("Failed to init local asset repo", log.ErrAttr(err))

				return err
			}

			assets := asset.NewAssetUsecase(assetRepository)
			serversUC := servers.NewServersUsecase(servers.NewServersRepository(dbConn))
			demos := demo.NewDemoUsecase(domain.BucketDemo, demo.NewDemoRepository(dbConn), assets, configUsecase, serversUC)

			reportUsecase := report.NewReportUsecase(report.NewReportRepository(dbConn), notificationUsecase, configUsecase, personUsecase, demos)

			stateUsecase := state.NewStateUsecase(eventBroadcaster,
				state.NewStateRepository(state.NewCollector(serversUC)), configUsecase, serversUC)

			banUsecase := ban.NewBanSteamUsecase(ban.NewBanSteamRepository(dbConn, personUsecase, networkUsecase), personUsecase, configUsecase, notificationUsecase, reportUsecase, stateUsecase)

			banGroupRepo := steamgroup.NewSteamGroupRepository(dbConn)
			banGroupUsecase := steamgroup.NewBanGroupUsecase(banGroupRepo, personUsecase, notificationUsecase, configUsecase)

			blocklistUsecase := blocklist.NewBlocklistUsecase(blocklist.NewBlocklistRepository(dbConn), banUsecase, banGroupUsecase)

			go func() {
				if err := stateUsecase.Start(ctx); err != nil {
					slog.Error("Failed to start state tracker", log.ErrAttr(err))
				}
			}()

			banASNUsecase := ban.NewBanASNUsecase(ban.NewBanASNRepository(dbConn), notificationUsecase, networkUsecase, configUsecase, personUsecase)
			banNetUsecase := ban.NewBanNetUsecase(ban.NewBanNetRepository(dbConn), personUsecase, configUsecase, notificationUsecase, stateUsecase)

			discordOAuthUsecase := discord.NewDiscordOAuthUsecase(discord.NewDiscordOAuthRepository(dbConn), configUsecase)

			appeals := appeal.NewAppealUsecase(appeal.NewAppealRepository(dbConn), banUsecase, personUsecase, notificationUsecase, configUsecase)

			matchRepo := match.NewMatchRepository(eventBroadcaster, dbConn, personUsecase, serversUC, notificationUsecase, stateUsecase, weaponsMap)
			go matchRepo.Start(ctx)

			matchUsecase := match.NewMatchUsecase(matchRepo, stateUsecase, serversUC, notificationUsecase)

			if errWeapons := matchUsecase.LoadWeapons(ctx, weaponsMap); errWeapons != nil {
				slog.Error("Failed to import weapons", log.ErrAttr(errWeapons))
			}

			chatRepository := chat.NewChatRepository(dbConn, personUsecase, wordFilterUsecase, matchUsecase, eventBroadcaster)
			go chatRepository.Start(ctx)

			chatUsecase := chat.NewChatUsecase(configUsecase, chatRepository, wordFilterUsecase, stateUsecase, banUsecase, personUsecase, notificationUsecase)
			go chatUsecase.Start(ctx)

			forumUsecase := forum.NewForumUsecase(forum.NewForumRepository(dbConn), notificationUsecase)
			go forumUsecase.Start(ctx)

			metricsUsecase := metrics.NewMetricsUsecase(eventBroadcaster)
			go metricsUsecase.Start(ctx)

			newsUsecase := news.NewNewsUsecase(news.NewNewsRepository(dbConn))
			patreonUsecase := patreon.NewPatreonUsecase(patreon.NewPatreonRepository(dbConn), configUsecase)
			srcdsUsecase := srcds.NewSrcdsUsecase(srcds.NewRepository(dbConn), configUsecase, serversUC, personUsecase, reportUsecase, notificationUsecase, banUsecase)
			wikiUsecase := wiki.NewWikiUsecase(wiki.NewWikiRepository(dbConn))
			authRepo := auth.NewAuthRepository(dbConn)
			authUsecase := auth.NewAuthUsecase(authRepo, configUsecase, personUsecase, banUsecase, serversUC)

			voteUsecase := votes.NewVoteUsecase(votes.NewVoteRepository(dbConn), personUsecase, matchUsecase, notificationUsecase, configUsecase, eventBroadcaster)
			go voteUsecase.Start(ctx)

			contestUsecase := contest.NewContestUsecase(contest.NewContestRepository(dbConn))

			if err := firstTimeSetup(ctx, personUsecase, newsUsecase, wikiUsecase, conf); err != nil {
				slog.Error("Failed to run first time setup", log.ErrAttr(err))

				return err
			}

			// start workers
			if conf.General.Mode == domain.ReleaseMode {
				gin.SetMode(gin.ReleaseMode)
			} else {
				gin.SetMode(gin.DebugMode)
			}

			router, err := httphelper.CreateRouter(conf, app.Version())
			if err != nil {
				slog.Error("Could not setup router", log.ErrAttr(err))

				return err
			}

			discordHandler := discord.NewDiscordHandler(discordUsecase, personUsecase, banUsecase,
				stateUsecase, serversUC, configUsecase, networkUsecase, wordFilterUsecase, matchUsecase, banNetUsecase, banASNUsecase)
			discordHandler.Start(ctx)

			appeal.NewAppealHandler(router, appeals, authUsecase)
			auth.NewAuthHandler(router, authUsecase, configUsecase, personUsecase)
			ban.NewBanHandler(router, banUsecase, discordUsecase, personUsecase, configUsecase, authUsecase)
			ban.NewBanNetHandler(router, banNetUsecase, authUsecase)
			ban.NewBanASNHandler(router, banASNUsecase, authUsecase)
			config.NewConfigHandler(router, configUsecase, authUsecase, app.Version())
			discord.NewDiscordOAuthHandler(router, authUsecase, configUsecase, personUsecase, discordOAuthUsecase)
			steamgroup.NewSteamgroupHandler(router, banGroupUsecase, authUsecase)
			blocklist.NewBlocklistHandler(router, blocklistUsecase, networkUsecase, authUsecase)
			chat.NewChatHandler(router, chatUsecase, authUsecase)
			contest.NewContestHandler(router, contestUsecase, configUsecase, assets, authUsecase)
			demo.NewDemoHandler(router, demos)
			forum.NewForumHandler(router, forumUsecase, authUsecase)
			match.NewMatchHandler(ctx, router, matchUsecase, serversUC, authUsecase, configUsecase)
			asset.NewAssetHandler(router, configUsecase, assets, authUsecase)
			metrics.NewMetricsHandler(router)
			network.NewNetworkHandler(router, networkUsecase, authUsecase)
			news.NewNewsHandler(router, newsUsecase, notificationUsecase, authUsecase)
			notification.NewNotificationHandler(router, notificationUsecase, authUsecase)
			patreon.NewPatreonHandler(router, patreonUsecase, authUsecase, configUsecase)
			person.NewPersonHandler(router, configUsecase, personUsecase, authUsecase)
			report.NewReportHandler(router, reportUsecase, authUsecase, notificationUsecase)
			servers.NewServerHandler(router, serversUC, stateUsecase, authUsecase, personUsecase)
			srcds.NewSRCDSHandler(router, srcdsUsecase, serversUC, personUsecase, assets,
				reportUsecase, banUsecase, networkUsecase, banGroupUsecase, demos, authUsecase, banASNUsecase, banNetUsecase,
				configUsecase, notificationUsecase, stateUsecase, blocklistUsecase)
			votes.NewVoteHandler(router, voteUsecase, authUsecase)
			wiki.NewWIkiHandler(router, wikiUsecase, authUsecase)
			wordfilter.NewWordFilterHandler(router, configUsecase, wordFilterUsecase, chatUsecase, authUsecase)

			if conf.Debug.AddRCONLogAddress != "" {
				go stateUsecase.LogAddressAdd(ctx, conf.Debug.AddRCONLogAddress)
			}

			// River Queue
			workers := createQueueWorkers(
				personUsecase,
				notificationUsecase,
				discordUsecase,
				authRepo,
				patreonUsecase,
				reportUsecase,
				discordOAuthUsecase)

			memberships := steamgroup.NewMemberships(banGroupRepo)
			banExpirations := ban.NewExpirationMonitor(banUsecase, banNetUsecase, banASNUsecase, personUsecase, notificationUsecase, configUsecase)

			go func() {
				go memberships.Update(ctx)
				go banExpirations.Update(ctx)
				go blocklistUsecase.Sync(ctx)
				go demos.Cleanup(ctx)

				membershipsTicker := time.NewTicker(12 * time.Hour)
				expirationsTicker := time.NewTicker(60 * time.Second)
				reportIntoTicker := time.NewTicker(24 * time.Hour)
				blocklistTicker := time.NewTicker(6 * time.Hour)
				demoTicker := time.NewTicker(5 * time.Minute)

				select {
				case <-ctx.Done():
					return
				case <-membershipsTicker.C:
					go memberships.Update(ctx)
				case <-expirationsTicker.C:
					go banExpirations.Update(ctx)
				case <-reportIntoTicker.C:
					go func() {
						if errMeta := reportUsecase.GenerateMetaStats(ctx); errMeta != nil {
							slog.Error("Failed to generate meta stats", log.ErrAttr(errMeta))
						}
					}()
				case <-blocklistTicker.C:
					go blocklistUsecase.Sync(ctx)
				case <-demoTicker.C:
					go demos.Cleanup(ctx)
				}
			}()

			periodicJons := createPeriodicJobs()
			queueClient, errClient := queue.Client(dbConn.Pool(), workers, periodicJons)
			if errClient != nil {
				slog.Error("Failed to setup job queue", log.ErrAttr(errClient))

				return errClient
			}

			if errClientStart := queueClient.Start(ctx); errClientStart != nil {
				slog.Error("Failed to start job client", log.ErrAttr(errClientStart))

				return errors.Join(errClientStart, queue.ErrStartQueue)
			}

			notificationUsecase.SetQueueClient(queueClient)

			httpServer := httphelper.NewHTTPServer(conf.Addr(), router)

			demoDownloader := demo.NewDownloader(configUsecase, dbConn, serversUC, assets, demos)
			go demoDownloader.Start(ctx)

			go func() {
				<-ctx.Done()

				slog.Info("Shutting down HTTP service")

				shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				if errShutdown := httpServer.Shutdown(shutdownCtx); errShutdown != nil { //nolint:contextcheck
					slog.Error("Error shutting down http service", log.ErrAttr(errShutdown))
				}
			}()

			errServe := httpServer.ListenAndServe()
			if errServe != nil && !errors.Is(errServe, http.ErrServerClosed) {
				slog.Error("HTTP server returned error", log.ErrAttr(errServe))
			}

			<-ctx.Done()

			slog.Info("Exiting...")

			return nil
		},
	}
}
