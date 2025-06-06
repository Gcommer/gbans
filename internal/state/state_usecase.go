package state

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/leighmacdonald/gbans/internal/domain"
	"github.com/leighmacdonald/gbans/pkg/fp"
	"github.com/leighmacdonald/gbans/pkg/log"
	"github.com/leighmacdonald/gbans/pkg/logparse"
	"github.com/leighmacdonald/steamid/v4/steamid"
	"github.com/ryanuber/go-glob"
	"golang.org/x/sync/errgroup"
)

type stateUsecase struct {
	state       domain.StateRepository
	config      domain.ConfigUsecase
	servers     domain.ServersUsecase
	logListener *logparse.UDPLogListener
	logFileChan chan domain.LogFilePayload
	broadcaster *fp.Broadcaster[logparse.EventType, logparse.ServerEvent]
}

// NewStateUsecase created a interface to interact with server state and exec rcon commands.
func NewStateUsecase(broadcaster *fp.Broadcaster[logparse.EventType, logparse.ServerEvent],
	repository domain.StateRepository, config domain.ConfigUsecase, servers domain.ServersUsecase,
) domain.StateUsecase {
	return &stateUsecase{
		state:       repository,
		config:      config,
		broadcaster: broadcaster,
		servers:     servers,
		logFileChan: make(chan domain.LogFilePayload),
	}
}

func (s *stateUsecase) Start(ctx context.Context) error {
	conf := s.config.Config()

	logSrc, errLogSrc := logparse.NewUDPLogListener(conf.General.SrcdsLogAddr,
		func(_ logparse.EventType, event logparse.ServerEvent) {
			s.broadcaster.Emit(event.EventType, event)
		})

	if errLogSrc != nil {
		return errLogSrc
	}

	s.logListener = logSrc

	go s.state.Start(ctx)

	// TODO run on server Config changes
	s.updateSrcdsLogServers(ctx)

	s.logListener.Start(ctx)

	return nil
}

func (s *stateUsecase) updateSrcdsLogServers(ctx context.Context) {
	newSecrets := map[int]logparse.ServerIDMap{}
	newServers := map[netip.Addr]bool{}
	serversCtx, cancelServers := context.WithTimeout(ctx, time.Second*5)

	defer cancelServers()

	servers, _, errServers := s.servers.Servers(serversCtx, domain.ServerQueryFilter{
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

		if ip, errIP := server.IP(ctx); errIP == nil {
			if addr, addrOk := netip.AddrFromSlice(ip); addrOk {
				newServers[addr] = true
			} else {
				slog.Error("Failed to convert ip", slog.String("address", ip.String()))
			}
		} else {
			slog.Error("Failed to resolve server ip", log.ErrAttr(errIP))
		}

		if internalIP, errIP := server.IPInternal(ctx); errIP == nil {
			if addr, addrOk := netip.AddrFromSlice(internalIP); addrOk {
				newServers[addr] = true
			} else {
				slog.Error("Failed to convert internal ip", slog.String("address", internalIP.String()))
			}
		} else {
			slog.Error("Failed to resolve internal server ip", log.ErrAttr(errIP))
		}
	}

	s.logListener.SetSecrets(newSecrets)
	s.logListener.SetServers(newServers)
}

func (s *stateUsecase) Current() []domain.ServerState {
	return s.state.Current()
}

func (s *stateUsecase) FindByCIDR(cidr *net.IPNet) []domain.PlayerServerInfo {
	return s.Find("", steamid.SteamID{}, nil, cidr)
}

func (s *stateUsecase) FindByIP(addr net.IP) []domain.PlayerServerInfo {
	return s.Find("", steamid.SteamID{}, addr, nil)
}

func (s *stateUsecase) FindByName(name string) []domain.PlayerServerInfo {
	return s.Find(name, steamid.SteamID{}, nil, nil)
}

func (s *stateUsecase) FindBySteamID(steamID steamid.SteamID) []domain.PlayerServerInfo {
	return s.Find("", steamID, nil, nil)
}

func (s *stateUsecase) Update(serverID int, update domain.PartialStateUpdate) error {
	return s.state.Update(serverID, update)
}

// Find searches the current server state for players matching at least one of the provided criteria.
func (s *stateUsecase) Find(name string, steamID steamid.SteamID, addr net.IP, cidr *net.IPNet) []domain.PlayerServerInfo {
	var found []domain.PlayerServerInfo

	current := s.state.Current()

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
	for _, server := range s.state.Current() {
		_, exists := serverMap[server.Region]
		if !exists {
			serverMap[server.Region] = []domain.ServerState{}
		}

		serverMap[server.Region] = append(serverMap[server.Region], server)
	}

	return serverMap
}

func (s *stateUsecase) ByServerID(serverID int) (domain.ServerState, bool) {
	for _, server := range s.state.Current() {
		if server.ServerID == serverID {
			return server, true
		}
	}

	return domain.ServerState{}, false
}

func (s *stateUsecase) ByName(name string, wildcardOk bool) []domain.ServerState {
	var servers []domain.ServerState

	current := s.state.Current()

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

func (s *stateUsecase) OnFindExec(ctx context.Context, name string, steamID steamid.SteamID, ip net.IP, cidr *net.IPNet, onFoundCmd func(info domain.PlayerServerInfo) string) error {
	currentState := s.state.Current()
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

	for _, server := range s.state.Configs() {
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
	return s.state.ExecRaw(ctx, addr, password, cmd)
}

func (s *stateUsecase) LogAddressAdd(ctx context.Context, logAddress string) {
	time.Sleep(20 * time.Second)
	slog.Info("Enabling log forwarding", slog.String("host", logAddress))
	s.Broadcast(ctx, nil, "logaddress_add "+logAddress)
}

func (s *stateUsecase) LogAddressDel(ctx context.Context, logAddress string) {
	slog.Info("Disabling log forwarding for host", slog.String("host", logAddress))
	s.Broadcast(ctx, nil, "logaddress_add "+logAddress)
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

	configs := s.state.Configs()

	if len(serverIDs) == 0 {
		for _, conf := range configs {
			serverIDs = append(serverIDs, conf.ServerID)
		}
	}

	resultChan := make(chan broadcastResult)

	for _, serverID := range serverIDs {
		sid := serverID

		errGroup.Go(func() error {
			serverConf, errServerConf := s.state.GetServer(sid)
			if errServerConf != nil {
				return errServerConf
			}

			resp, errExec := s.state.ExecRaw(egCtx, serverConf.Addr(), serverConf.RconPassword, cmd)
			if errExec != nil {
				if errors.Is(errExec, context.Canceled) {
					return nil
				}

				slog.Error("Failed to exec server command", slog.String("name", serverConf.DefaultHostname),
					slog.Int("server_id", sid), log.ErrAttr(errExec))

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
func (s *stateUsecase) Kick(ctx context.Context, target steamid.SteamID, reason domain.Reason) error {
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

// KickPlayerID will kick the steam id from whatever server it is connected to.
func (s *stateUsecase) KickPlayerID(ctx context.Context, targetPlayerID int, targetServerID int, reason domain.Reason) error {
	_, err := s.ExecServer(ctx, targetServerID, fmt.Sprintf("sm_kick #%d %s", targetPlayerID, reason.String()))

	return err
}

// Silence will gag & mute a player.
func (s *stateUsecase) Silence(ctx context.Context, target steamid.SteamID, reason domain.Reason,
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

		return fmt.Sprintf(`sm_silence "#%s" %s`, info.Player.SID.Steam(false), reason.String())
	}); errExec != nil {
		return errors.Join(errExec, fmt.Errorf("%w: sm_silence", domain.ErrCommandFailed))
	}

	return nil
}

// Say is used to send a message to the server via sm_say.
func (s *stateUsecase) Say(ctx context.Context, serverID int, message string) error {
	_, errExec := s.ExecServer(ctx, serverID, `sm_say `+message)

	return errors.Join(errExec, fmt.Errorf("%w: sm_say", domain.ErrCommandFailed))
}

// CSay is used to send a centered message to the server via sm_csay.
func (s *stateUsecase) CSay(ctx context.Context, serverID int, message string) error {
	_, errExec := s.ExecServer(ctx, serverID, `sm_csay `+message)

	return errors.Join(errExec, fmt.Errorf("%w: sm_csay", domain.ErrCommandFailed))
}

// PSay is used to send a private message to a player.
func (s *stateUsecase) PSay(ctx context.Context, target steamid.SteamID, message string) error {
	if !target.Valid() {
		return domain.ErrInvalidTargetSID
	}

	if errExec := s.OnFindExec(ctx, "", target, nil, nil, func(_ domain.PlayerServerInfo) string {
		return fmt.Sprintf(`sm_psay "#%s" "%s"`, target.Steam(false), message)
	}); errExec != nil {
		return errors.Join(errExec, fmt.Errorf("%w: sm_psay", domain.ErrCommandFailed))
	}

	return nil
}
