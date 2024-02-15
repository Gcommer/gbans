package state

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/leighmacdonald/gbans/internal/domain"
	"github.com/leighmacdonald/gbans/pkg/fp"
	"github.com/leighmacdonald/gbans/pkg/log"
	"github.com/leighmacdonald/gbans/pkg/logparse"
	"github.com/leighmacdonald/steamid/v3/steamid"
	"github.com/ryanuber/go-glob"
	"golang.org/x/sync/errgroup"
)

type stateUsecase struct {
	stateRepository domain.StateRepository
	configUsecase   domain.ConfigUsecase
	serversUsecase  domain.ServersUsecase
	logListener     *logparse.UDPLogListener
	logFileChan     chan domain.LogFilePayload
	broadcaster     *fp.Broadcaster[logparse.EventType, logparse.ServerEvent]
}

// NewStateUsecase created a interface to interact with server state and exec rcon commands
// TODO ensure started.
func NewStateUsecase(broadcaster *fp.Broadcaster[logparse.EventType, logparse.ServerEvent],
	repository domain.StateRepository, configUsecase domain.ConfigUsecase, serversUsecase domain.ServersUsecase,
) domain.StateUsecase {
	return &stateUsecase{
		stateRepository: repository,
		configUsecase:   configUsecase,
		broadcaster:     broadcaster,
		serversUsecase:  serversUsecase,
		logFileChan:     make(chan domain.LogFilePayload),
	}
}

func (s *stateUsecase) Start(ctx context.Context) error {
	conf := s.configUsecase.Config()

	logSrc, errLogSrc := logparse.NewUDPLogListener(conf.Log.SrcdsLogAddr,
		func(eventType logparse.EventType, event logparse.ServerEvent) {
			s.broadcaster.Emit(event.EventType, event)
		})

	if errLogSrc != nil {
		return errLogSrc
	}

	s.logListener = logSrc

	s.stateRepository.Start(ctx)

	s.updateSrcdsLogSecrets(ctx)

	go s.logReader(ctx, false)

	// TODO run on server Config changes
	s.updateSrcdsLogSecrets(ctx)

	s.logListener.Start(ctx)

	return nil
}

// logReader is the fan-out orchestrator for game log events
// Registering receivers can be accomplished with app.eb.Broadcaster.
func (s *stateUsecase) logReader(ctx context.Context, writeUnhandled bool) {
	var (
		logger = slog.Default().WithGroup("logReader")
		file   *os.File
	)

	if writeUnhandled {
		var errCreateFile error
		file, errCreateFile = os.Create("./unhandled_messages.log")

		if errCreateFile != nil {
			logger.Error("Failed to open debug message log", log.ErrAttr(errCreateFile))
		} else {
			defer func() {
				if errClose := file.Close(); errClose != nil {
					logger.Error("Failed to close unhandled_messages.log", log.ErrAttr(errClose))
				}
			}()
		}
	}

	parser := logparse.NewLogParser()

	// playerStateCache := newPlayerCache(app.logger)
	for {
		select {
		case logFile := <-s.logFileChan:
			emitted := 0
			failed := 0
			unknown := 0
			ignored := 0

			for _, logLine := range logFile.Lines {
				parseResult, errParse := parser.Parse(logLine)
				if errParse != nil {
					continue
				}

				newServerEvent := logparse.ServerEvent{
					ServerName: logFile.ServerName,
					ServerID:   logFile.ServerID,
					Results:    parseResult,
				}

				if newServerEvent.EventType == logparse.IgnoredMsg {
					ignored++

					continue
				} else if newServerEvent.EventType == logparse.UnknownMsg {
					unknown++

					if writeUnhandled {
						if _, errWrite := file.WriteString(logLine + "\n"); errWrite != nil {
							logger.Error("Failed to write debug log", log.ErrAttr(errWrite))
						}
					}
				}

				s.broadcaster.Emit(newServerEvent.EventType, newServerEvent)

				emitted++
			}

			logger.Debug("Completed emitting logfile events",
				slog.Int("ok", emitted), slog.Int("failed", failed),
				slog.Int("unknown", unknown), slog.Int("ignored", ignored))
		case <-ctx.Done():
			logger.Debug("logReader shutting down")

			return
		}
	}
}

func (s *stateUsecase) updateSrcdsLogSecrets(ctx context.Context) {
	newSecrets := map[int]logparse.ServerIDMap{}
	serversCtx, cancelServers := context.WithTimeout(ctx, time.Second*5)

	defer cancelServers()

	servers, _, errServers := s.serversUsecase.GetServers(serversCtx, domain.ServerQueryFilter{
		IncludeDisabled: false,
		QueryFilter:     domain.QueryFilter{Deleted: false},
	})
	if errServers != nil {
		slog.Error("Failed to update srcds log secrets", log.ErrAttr(errServers))

		return
	}

	for _, server := range servers {
		newSecrets[server.LogSecret] = logparse.ServerIDMap{
			ServerID:   server.ServerID,
			ServerName: server.ShortName,
		}
	}

	s.logListener.SetSecrets(newSecrets)
}

func (s *stateUsecase) Current() []domain.ServerState {
	return s.stateRepository.Current()
}

func (s *stateUsecase) FindByCIDR(cidr *net.IPNet) []domain.PlayerServerInfo {
	return s.Find("", "", nil, cidr)
}

func (s *stateUsecase) FindByIP(addr net.IP) []domain.PlayerServerInfo {
	return s.Find("", "", addr, nil)
}

func (s *stateUsecase) FindByName(name string) []domain.PlayerServerInfo {
	return s.Find(name, "", nil, nil)
}

func (s *stateUsecase) FindBySteamID(steamID steamid.SID64) []domain.PlayerServerInfo {
	return s.Find("", steamID, nil, nil)
}

func (s *stateUsecase) Update(serverID int, update domain.PartialStateUpdate) error {
	return s.stateRepository.Update(serverID, update)
}

func (s *stateUsecase) Find(name string, steamID steamid.SID64, addr net.IP, cidr *net.IPNet) []domain.PlayerServerInfo {
	var found []domain.PlayerServerInfo

	current := s.stateRepository.Current()

	for server := range current {
		for _, player := range current[server].Players {
			matched := false
			if steamID.Valid() && player.SID == steamID {
				matched = true
			}

			if name != "" {
				queryName := name
				if !strings.HasPrefix(queryName, "*") {
					queryName = "*" + queryName
				}

				if !strings.HasSuffix(queryName, "*") {
					queryName += "*"
				}

				m := glob.Glob(strings.ToLower(queryName), strings.ToLower(player.Name))
				if m {
					matched = true
				}
			}

			if addr != nil && addr.Equal(player.IP) {
				matched = true
			}

			if cidr != nil && cidr.Contains(player.IP) {
				matched = true
			}

			if matched {
				found = append(found, domain.PlayerServerInfo{Player: player, ServerID: current[server].ServerID})
			}
		}
	}

	return found
}

func (s *stateUsecase) SortRegion() map[string][]domain.ServerState {
	serverMap := map[string][]domain.ServerState{}
	for _, server := range s.stateRepository.Current() {
		_, exists := serverMap[server.Region]
		if !exists {
			serverMap[server.Region] = []domain.ServerState{}
		}

		serverMap[server.Region] = append(serverMap[server.Region], server)
	}

	return serverMap
}

func (s *stateUsecase) ByServerID(serverID int) (domain.ServerState, bool) {
	for _, server := range s.stateRepository.Current() {
		if server.ServerID == serverID {
			return server, true
		}
	}

	return domain.ServerState{}, false
}

func (s *stateUsecase) ByName(name string, wildcardOk bool) []domain.ServerState {
	var servers []domain.ServerState

	current := s.stateRepository.Current()

	if name == "*" && wildcardOk {
		servers = append(servers, current...)
	} else {
		if !strings.HasPrefix(name, "*") {
			name = "*" + name
		}

		if !strings.HasSuffix(name, "*") {
			name += "*"
		}

		for _, server := range current {
			if glob.Glob(strings.ToLower(name), strings.ToLower(server.NameShort)) ||
				strings.EqualFold(server.NameShort, name) {
				servers = append(servers, server)

				break
			}
		}
	}

	return servers
}

func (s *stateUsecase) ServerIDsByName(name string, wildcardOk bool) []int {
	var servers []int //nolint:prealloc
	for _, server := range s.ByName(name, wildcardOk) {
		servers = append(servers, server.ServerID)
	}

	return servers
}

func (s *stateUsecase) OnFindExec(ctx context.Context, name string, steamID steamid.SID64, ip net.IP, cidr *net.IPNet, onFoundCmd func(info domain.PlayerServerInfo) string) error {
	currentState := s.stateRepository.Current()
	players := s.Find(name, steamID, ip, cidr)

	if len(players) == 0 {
		return domain.ErrPlayerNotFound
	}

	var err error

	for _, player := range players {
		for _, server := range currentState {
			if player.ServerID == server.ServerID {
				_, errRcon := s.ExecServer(ctx, server.ServerID, onFoundCmd(player))
				if errRcon != nil {
					err = errors.Join(errRcon)
				}
			}
		}
	}

	return err
}

func (s *stateUsecase) ExecServer(ctx context.Context, serverID int, cmd string) (string, error) {
	var conf domain.ServerConfig

	for _, server := range s.stateRepository.Configs() {
		if server.ServerID == serverID {
			conf = server

			break
		}
	}

	if conf.ServerID == 0 {
		return "", domain.ErrUnknownServerID
	}

	return s.ExecRaw(ctx, conf.Addr(), conf.RconPassword, cmd)
}

func (s *stateUsecase) ExecRaw(ctx context.Context, addr string, password string, cmd string) (string, error) {
	return s.stateRepository.ExecRaw(ctx, addr, password, cmd)
}

func (s *stateUsecase) LogAddressAdd(ctx context.Context, logAddress string) {
	s.Broadcast(ctx, nil, fmt.Sprintf("logaddress_add %s", logAddress))
}

type broadcastResult struct {
	serverID int
	resp     string
}

// Broadcast sends out rcon commands to all provided servers. If no servers are provided it will default to broadcasting
// to every server.
func (s *stateUsecase) Broadcast(ctx context.Context, serverIDs []int, cmd string) map[int]string {
	results := map[int]string{}
	errGroup, egCtx := errgroup.WithContext(ctx)

	configs := s.stateRepository.Configs()

	if len(serverIDs) == 0 {
		for _, conf := range configs {
			serverIDs = append(serverIDs, conf.ServerID)
		}
	}

	resultChan := make(chan broadcastResult)

	for _, serverID := range serverIDs {
		sid := serverID

		errGroup.Go(func() error {
			serverConf, errServerConf := s.stateRepository.GetServer(sid)
			if errServerConf != nil {
				return errServerConf
			}

			resp, errExec := s.stateRepository.ExecRaw(egCtx, serverConf.Addr(), serverConf.RconPassword, cmd)
			if errExec != nil {
				slog.Error("Failed to exec server command", slog.Int("server_id", sid), log.ErrAttr(errExec))

				// Don't error out since we don't want a single servers potentially temporary issue to prevent the rest
				// from executing.
				return nil
			}

			resultChan <- broadcastResult{
				serverID: sid,
				resp:     resp,
			}

			return nil
		})
	}

	go func() {
		err := errGroup.Wait()
		if err != nil {
			slog.Error("Failed to broadcast command", log.ErrAttr(err))
		}

		close(resultChan)
	}()

	for result := range resultChan {
		results[result.serverID] = result.resp
	}

	return results
}

// Kick will kick the steam id from whatever server it is connected to.
func (s *stateUsecase) Kick(ctx context.Context, target steamid.SID64, reason domain.Reason) error {
	if !target.Valid() {
		return domain.ErrInvalidTargetSID
	}

	if errExec := s.OnFindExec(ctx, "", target, nil, nil, func(info domain.PlayerServerInfo) string {
		return fmt.Sprintf("sm_kick #%d %s", info.Player.UserID, reason.String())
	}); errExec != nil {
		return errors.Join(errExec, domain.ErrCommandFailed)
	}

	return nil
}

// Silence will gag & mute a player.
func (s *stateUsecase) Silence(ctx context.Context, target steamid.SID64, reason domain.Reason,
) error {
	if !target.Valid() {
		return domain.ErrInvalidTargetSID
	}

	var (
		users   []string
		usersMu = &sync.RWMutex{}
	)

	if errExec := s.OnFindExec(ctx, "", target, nil, nil, func(info domain.PlayerServerInfo) string {
		usersMu.Lock()
		users = append(users, info.Player.Name)
		usersMu.Unlock()

		return fmt.Sprintf(`sm_silence "#%s" %s`, steamid.SID64ToSID(info.Player.SID), reason.String())
	}); errExec != nil {
		return errors.Join(errExec, fmt.Errorf("%w: sm_silence", domain.ErrCommandFailed))
	}

	return nil
}

// Say is used to send a message to the server via sm_say.
func (s *stateUsecase) Say(ctx context.Context, serverID int, message string) error {
	_, errExec := s.ExecServer(ctx, serverID, fmt.Sprintf(`sm_say %s`, message))

	return errors.Join(errExec, fmt.Errorf("%w: sm_say", domain.ErrCommandFailed))
}

// CSay is used to send a centered message to the server via sm_csay.
func (s *stateUsecase) CSay(ctx context.Context, serverID int, message string) error {
	_, errExec := s.ExecServer(ctx, serverID, fmt.Sprintf(`sm_csay %s`, message))

	return errors.Join(errExec, fmt.Errorf("%w: sm_csay", domain.ErrCommandFailed))
}

// PSay is used to send a private message to a player.
func (s *stateUsecase) PSay(ctx context.Context, target steamid.SID64, message string) error {
	if !target.Valid() {
		return domain.ErrInvalidTargetSID
	}

	if errExec := s.OnFindExec(ctx, "", target, nil, nil, func(info domain.PlayerServerInfo) string {
		return fmt.Sprintf(`sm_psay "#%s" "%s"`, steamid.SID64ToSID(target), message)
	}); errExec != nil {
		return errors.Join(errExec, fmt.Errorf("%w: sm_psay", domain.ErrCommandFailed))
	}

	return nil
}