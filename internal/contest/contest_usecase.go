package contest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/leighmacdonald/gbans/internal/domain"
	"github.com/leighmacdonald/steamid/v4/steamid"
)

type contestUsecase struct {
	repository domain.ContestRepository
}

func NewContestUsecase(repository domain.ContestRepository) domain.ContestUsecase {
	return &contestUsecase{repository: repository}
}

func (c *contestUsecase) ContestSave(ctx context.Context, contest domain.Contest) (domain.Contest, error) {
	if contest.ContestID.IsNil() {
		newID, errID := uuid.NewV4()
		if errID != nil {
			return contest, errors.Join(errID, domain.ErrUUIDCreate)
		}

		contest.ContestID = newID
	}

	if errSave := c.repository.ContestSave(ctx, &contest); errSave != nil {
		return contest, errSave
	}

	slog.Info("Contest updated",
		slog.String("contest_id", contest.ContestID.String()),
		slog.String("title", contest.Title))

	return contest, nil
}

func (c *contestUsecase) ContestByID(ctx context.Context, contestID uuid.UUID, contest *domain.Contest) error {
	return c.repository.ContestByID(ctx, contestID, contest)
}

func (c *contestUsecase) ContestDelete(ctx context.Context, contestID uuid.UUID) error {
	if err := c.repository.ContestDelete(ctx, contestID); err != nil {
		return err
	}

	slog.Info("Contest deleted", slog.String("contest_id", contestID.String()))

	return nil
}

func (c *contestUsecase) ContestEntryDelete(ctx context.Context, contestEntryID uuid.UUID) error {
	return c.repository.ContestEntryDelete(ctx, contestEntryID)
}

func (c *contestUsecase) Contests(ctx context.Context, user domain.PersonInfo) ([]domain.Contest, error) {
	return c.repository.Contests(ctx, !user.HasPermission(domain.PModerator))
}

func (c *contestUsecase) ContestEntry(ctx context.Context, contestID uuid.UUID, entry *domain.ContestEntry) error {
	return c.repository.ContestEntry(ctx, contestID, entry)
}

func (c *contestUsecase) ContestEntrySave(ctx context.Context, entry domain.ContestEntry) error {
	return c.repository.ContestEntrySave(ctx, entry)
}

func (c *contestUsecase) ContestEntries(ctx context.Context, contestID uuid.UUID) ([]*domain.ContestEntry, error) {
	return c.repository.ContestEntries(ctx, contestID)
}

func (c *contestUsecase) ContestEntryVoteGet(ctx context.Context, contestEntryID uuid.UUID, steamID steamid.SteamID, record *domain.ContentVoteRecord) error {
	return c.repository.ContestEntryVoteGet(ctx, contestEntryID, steamID, record)
}

func (c *contestUsecase) ContestEntryVote(ctx context.Context, contestID uuid.UUID, contestEntryID uuid.UUID, user domain.PersonInfo, vote bool) error {
	var contest domain.Contest
	if errContests := c.ContestByID(ctx, contestID, &contest); errContests != nil {
		return errContests
	}

	if !contest.Public && !user.HasPermission(domain.PModerator) {
		return domain.ErrPermissionDenied
	}

	if !contest.Voting || !contest.DownVotes && !vote {
		return domain.ErrBadRequest
	}

	if err := c.repository.ContestEntryVote(ctx, contestEntryID, user.GetSteamID(), vote); err != nil {
		return err
	}

	sid := user.GetSteamID()

	slog.Info("Entry vote registered", slog.String("contest_id", contest.ContestID.String()), slog.Bool("vote", vote), slog.String("steam_id", sid.String()))

	return nil
}

func (c *contestUsecase) ContestEntryVoteDelete(ctx context.Context, contestEntryVoteID int64) error {
	return c.repository.ContestEntryVoteDelete(ctx, contestEntryVoteID)
}

func (c *contestUsecase) ContestEntryVoteUpdate(ctx context.Context, contestEntryVoteID int64, newVote bool) error {
	return c.repository.ContestEntryVoteUpdate(ctx, contestEntryVoteID, newVote)
}
