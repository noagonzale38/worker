package messagequeue

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TicketsBot-cloud/common/sentry"
	"github.com/TicketsBot-cloud/database"
	"github.com/TicketsBot-cloud/gdl/rest"
	gdlUtils "github.com/TicketsBot-cloud/gdl/utils"
	workercontext "github.com/TicketsBot-cloud/worker"
	"github.com/TicketsBot-cloud/worker/bot/cache"
	cmdcontext "github.com/TicketsBot-cloud/worker/bot/command/context"
	"github.com/TicketsBot-cloud/worker/bot/constants"
	"github.com/TicketsBot-cloud/worker/bot/dbclient"
	"github.com/TicketsBot-cloud/worker/bot/logic"
	"github.com/TicketsBot-cloud/worker/bot/metrics/statsd"
	"github.com/TicketsBot-cloud/worker/bot/redis"
	"github.com/TicketsBot-cloud/worker/bot/utils"
	"go.uber.org/zap"
)

const AutoCloseReason = "Automatically closed due to inactivity"
const DefaultAutoCloseWarningMessage = "%user%, this ticket will automatically close soon if there is no response."

type autoCloseAction string

const (
	autoCloseActionClose   autoCloseAction = "close"
	autoCloseActionWarning autoCloseAction = "warning"
	autoCloseRedisChannel  string          = "tickets:autoclose"
)

type autoCloseTicket struct {
	GuildId       uint64          `json:"guild_id"`
	TicketId      int             `json:"ticket_id"`
	LastMessageId *uint64         `json:"last_message_id"`
	Action        autoCloseAction `json:"action,omitempty"`
}

func ListenAutoClose(logger *zap.Logger) {
	ch := make(chan autoCloseTicket)
	go listenAutoClose(ch)

	for acTicket := range ch {
		statsd.Client.IncrementKey(statsd.AutoClose)

		acTicket := acTicket
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), constants.TimeoutCloseTicket)
			defer cancel()

			logger.Debug("Processing autoclose event",
				zap.Int("ticket_id", acTicket.TicketId),
				zap.Uint64("guild_id", acTicket.GuildId),
			)

			// get ticket
			ticket, err := dbclient.Client.Tickets.Get(ctx, acTicket.TicketId, acTicket.GuildId)
			if err != nil {
				logger.Error("Failed to fetch ticket for autoclose",
					zap.Int("ticket_id", acTicket.TicketId),
					zap.Uint64("guild_id", acTicket.GuildId),
					zap.Error(err),
				)
				sentry.Error(err)
				return
			}

			// get worker
			worker, err := buildContext(ctx, ticket, cache.Client)
			if err != nil {
				logger.Error("Failed to build worker context for autoclose",
					zap.Int("ticket_id", acTicket.TicketId),
					zap.Uint64("guild_id", acTicket.GuildId),
					zap.Error(err),
				)
				sentry.Error(err)
				return
			}

			// query already checks, but just to be sure
			if ticket.ChannelId == nil {
				logger.Warn("Ticket channel ID is nil for autoclose",
					zap.Int("ticket_id", acTicket.TicketId),
					zap.Uint64("guild_id", acTicket.GuildId),
				)
				return
			}

			stale, err := isStaleAutoCloseTicket(ctx, ticket, acTicket.LastMessageId)
			if err != nil {
				logger.Error("Failed to check ticket activity for autoclose",
					zap.Int("ticket_id", acTicket.TicketId),
					zap.Uint64("guild_id", acTicket.GuildId),
					zap.Error(err),
				)
				sentry.Error(err)
				return
			}

			if stale {
				logger.Info("Skipping stale autoclose event",
					zap.Int("ticket_id", acTicket.TicketId),
					zap.Uint64("guild_id", acTicket.GuildId),
					zap.String("action", string(acTicket.Action)),
				)
				return
			}

			// get premium status
			premiumTier, err := utils.PremiumClient.GetTierByGuildId(ctx, ticket.GuildId, true, worker.Token, worker.RateLimiter)
			if err != nil {
				logger.Error("Failed to get premium tier for autoclose",
					zap.Int("ticket_id", acTicket.TicketId),
					zap.Uint64("guild_id", acTicket.GuildId),
					zap.Error(err),
				)
				sentry.Error(err)
				return
			}

			switch acTicket.Action {
			case autoCloseActionWarning:
				if err := sendAutoCloseWarning(ctx, worker, ticket, acTicket.LastMessageId); err != nil {
					logger.Error("Failed to send autoclose warning",
						zap.Int("ticket_id", acTicket.TicketId),
						zap.Uint64("guild_id", acTicket.GuildId),
						zap.Error(err),
					)
					sentry.Error(err)
					return
				}
			case "", autoCloseActionClose:
				cc := cmdcontext.NewAutoCloseContext(ctx, worker, ticket.GuildId, *ticket.ChannelId, worker.BotId, premiumTier)
				logic.CloseTicket(ctx, cc, gdlUtils.StrPtr(AutoCloseReason), true)
			default:
				logger.Warn("Skipping unknown autoclose action",
					zap.Int("ticket_id", acTicket.TicketId),
					zap.Uint64("guild_id", acTicket.GuildId),
					zap.String("action", string(acTicket.Action)),
				)
				return
			}

			logger.Info("Successfully processed autoclose event",
				zap.Int("ticket_id", acTicket.TicketId),
				zap.Uint64("guild_id", acTicket.GuildId),
				zap.String("action", string(acTicket.Action)),
			)
		}()
	}
}

func listenAutoClose(ch chan autoCloseTicket) {
	for {
		data, err := redis.Client.BLPop(context.Background(), 0, autoCloseRedisChannel).Result()
		if err != nil || len(data) < 2 {
			continue
		}

		var ticket autoCloseTicket
		if err := json.Unmarshal([]byte(data[1]), &ticket); err != nil {
			continue
		}

		ch <- ticket
	}
}

func isStaleAutoCloseTicket(ctx context.Context, ticket database.Ticket, payloadLastMessageId *uint64) (bool, error) {
	lastMessage, err := dbclient.Client.TicketLastMessage.Get(ctx, ticket.GuildId, ticket.Id)
	if err != nil {
		return false, err
	}

	return !sameUint64Ptr(payloadLastMessageId, lastMessage.LastMessageId), nil
}

func sendAutoCloseWarning(ctx context.Context, worker *workercontext.Context, ticket database.Ticket, payloadLastMessageId *uint64) error {
	settings, err := dbclient.Client.AutoClose.Get(ctx, ticket.GuildId)
	if err != nil {
		return err
	}

	content := strings.TrimSpace(DefaultAutoCloseWarningMessage)
	if settings.WarningMessage != nil && strings.TrimSpace(*settings.WarningMessage) != "" {
		content = strings.TrimSpace(*settings.WarningMessage)
	}

	content = formatAutoCloseWarningMessage(content, ticket)
	if len(content) > 2000 {
		content = utils.StringMax(content, 2000)
	}

	if _, err := worker.CreateMessageComplex(*ticket.ChannelId, rest.CreateMessageData{
		Content: content,
	}); err != nil {
		return err
	}

	return dbclient.Client.AutoCloseWarnings.MarkSent(ctx, ticket.GuildId, ticket.Id, payloadLastMessageId)
}

func formatAutoCloseWarningMessage(content string, ticket database.Ticket) string {
	now := time.Now().Unix()
	channel := ""
	if ticket.ChannelId != nil {
		channel = fmt.Sprintf("<#%d>", *ticket.ChannelId)
	}

	return strings.NewReplacer(
		"%user_id%", strconv.FormatUint(ticket.UserId, 10),
		"%user%", fmt.Sprintf("<@%d>", ticket.UserId),
		"%ticket_id%", strconv.Itoa(ticket.Id),
		"%guild_id%", strconv.FormatUint(ticket.GuildId, 10),
		"%channel%", channel,
		"%time%", fmt.Sprintf("<t:%d:t>", now),
		"%date%", fmt.Sprintf("<t:%d:d>", now),
		"%datetime%", fmt.Sprintf("<t:%d:f>", now),
	).Replace(content)
}

func sameUint64Ptr(a, b *uint64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}

	return *a == *b
}
