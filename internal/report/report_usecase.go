package report

import (
	"context"
	"log/slog"
	"time"

	"github.com/leighmacdonald/gbans/internal/discord"
	"github.com/leighmacdonald/gbans/internal/domain"
	"github.com/leighmacdonald/gbans/internal/httphelper"
	"github.com/leighmacdonald/gbans/pkg/log"
	"github.com/leighmacdonald/steamid/v3/steamid"
)

type reportUsecase struct {
	rr domain.ReportRepository
	du domain.DiscordUsecase
	cu domain.ConfigUsecase
	pu domain.PersonUsecase
}

func NewReportUsecase(repository domain.ReportRepository, discordUsecase domain.DiscordUsecase,
	configUsecase domain.ConfigUsecase, personUsecase domain.PersonUsecase,
) domain.ReportUsecase {
	return &reportUsecase{
		du: discordUsecase,
		rr: repository,
		cu: configUsecase,
		pu: personUsecase,
	}
}

func (r reportUsecase) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Hour * 24)
	updateChan := make(chan any)

	go func() {
		time.Sleep(time.Second * 5)
		updateChan <- true
	}()

	for {
		select {
		case <-ticker.C:
			updateChan <- true
		case <-updateChan:
			reports, _, errReports := r.GetReports(ctx, domain.ReportQueryFilter{
				QueryFilter: domain.QueryFilter{
					Limit: 0,
				},
			})
			if errReports != nil {
				slog.Error("failed to fetch reports for report metadata", log.ErrAttr(errReports))

				continue
			}

			var (
				now  = time.Now()
				meta domain.ReportMeta
			)

			for _, report := range reports {
				if report.ReportStatus == domain.ClosedWithAction || report.ReportStatus == domain.ClosedWithoutAction {
					meta.TotalClosed++

					continue
				}

				meta.TotalOpen++

				if report.ReportStatus == domain.NeedMoreInfo {
					meta.NeedInfo++
				} else {
					meta.Open++
				}

				switch {
				case now.Sub(report.CreatedOn) > time.Hour*24*7:
					meta.OpenWeek++
				case now.Sub(report.CreatedOn) > time.Hour*24*3:
					meta.Open3Days++
				case now.Sub(report.CreatedOn) > time.Hour*24:
					meta.Open1Day++
				default:
					meta.OpenNew++
				}
			}

			r.du.SendPayload(domain.ChannelMod, discord.ReportStatsMessage(meta, r.cu.ExtURLRaw("/admin/reports")))
		case <-ctx.Done():
			slog.Debug("showReportMeta shutting down")

			return
		}
	}
}

func (r reportUsecase) GetReportBySteamID(ctx context.Context, authorID steamid.SID64, steamID steamid.SID64) (domain.Report, error) {
	return r.rr.GetReportBySteamID(ctx, authorID, steamID)
}

func (r reportUsecase) GetReports(ctx context.Context, opts domain.ReportQueryFilter) ([]domain.Report, int64, error) {
	return r.rr.GetReports(ctx, opts)
}

func (r reportUsecase) GetReport(ctx context.Context, curUser domain.PersonInfo, reportID int64) (domain.Report, error) {
	report, err := r.rr.GetReport(ctx, reportID)
	if err != nil {
		return domain.Report{}, err
	}

	author, errAuthor := r.pu.GetPersonBySteamID(ctx, report.SourceID)
	if errAuthor != nil {
		return report, errAuthor
	}

	if !httphelper.HasPrivilege(curUser, steamid.Collection{author.SteamID}, domain.PModerator) {
		return report, domain.ErrPermissionDenied
	}

	return report, nil
}

func (r reportUsecase) GetReportMessages(ctx context.Context, reportID int64) ([]domain.ReportMessage, error) {
	return r.rr.GetReportMessages(ctx, reportID)
}

func (r reportUsecase) GetReportMessageByID(ctx context.Context, reportMessageID int64) (domain.ReportMessage, error) {
	return r.rr.GetReportMessageByID(ctx, reportMessageID)
}

func (r reportUsecase) DropReportMessage(ctx context.Context, message *domain.ReportMessage) error {
	return r.rr.DropReportMessage(ctx, message)
}

func (r reportUsecase) DropReport(ctx context.Context, report *domain.Report) error {
	return r.rr.DropReport(ctx, report)
}

func (r reportUsecase) SaveReport(ctx context.Context, report *domain.Report) error {
	return r.rr.SaveReport(ctx, report)
}

func (r reportUsecase) SaveReportMessage(ctx context.Context, message *domain.ReportMessage) error {
	return r.rr.SaveReportMessage(ctx, message)
}