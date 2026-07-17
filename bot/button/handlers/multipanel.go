package handlers

import (
	"github.com/TicketsBot-cloud/common/sentry"
	"github.com/TicketsBot-cloud/worker/bot/button/registry"
	"github.com/TicketsBot-cloud/worker/bot/button/registry/matcher"
	"github.com/TicketsBot-cloud/worker/bot/command/context"
	"github.com/TicketsBot-cloud/worker/bot/constants"
	"github.com/TicketsBot-cloud/worker/bot/dbclient"
	"github.com/TicketsBot-cloud/worker/bot/logic"
)

type MultiPanelHandler struct{}

func (h *MultiPanelHandler) Matcher() matcher.Matcher {
	return &matcher.SimpleMatcher{
		CustomId: "multipanel",
	}
}

func (h *MultiPanelHandler) Properties() registry.Properties {
	return registry.Properties{
		Flags:   registry.SumFlags(registry.GuildAllowed),
		Timeout: constants.TimeoutOpenTicket,
	}
}

func (h *MultiPanelHandler) Execute(ctx *context.SelectMenuContext) {
	if len(ctx.InteractionData.Values) == 0 {
		return
	}

	panelCustomId := ctx.InteractionData.Values[0]

	panel, ok, err := dbclient.Client.Panel.GetByCustomId(ctx, ctx.GuildId(), panelCustomId)
	if err != nil {
		sentry.Error(err) // TODO: Proper context
		return
	}

	if ok {
		// TODO: Log this
		if panel.GuildId != ctx.GuildId() {
			return
		}

		// Validate panel access
		canProceed, outOfHoursTitle, outOfHoursWarning, outOfHoursColour, err := logic.ValidatePanelAccess(ctx, panel)
		if err != nil {
			ctx.HandleError(err)
			return
		}

		if !canProceed {
			return
		}

		startPanelFlow(ctx, panel, outOfHoursTitle, outOfHoursWarning, outOfHoursColour)
	}
}
