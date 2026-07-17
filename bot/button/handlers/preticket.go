package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/TicketsBot-cloud/database"
	"github.com/TicketsBot-cloud/gdl/objects/channel/embed"
	"github.com/TicketsBot-cloud/gdl/objects/interaction/component"
	"github.com/TicketsBot-cloud/worker/bot/button"
	"github.com/TicketsBot-cloud/worker/bot/button/registry"
	"github.com/TicketsBot-cloud/worker/bot/button/registry/matcher"
	"github.com/TicketsBot-cloud/worker/bot/command"
	"github.com/TicketsBot-cloud/worker/bot/command/context"
	cmdregistry "github.com/TicketsBot-cloud/worker/bot/command/registry"
	"github.com/TicketsBot-cloud/worker/bot/constants"
	"github.com/TicketsBot-cloud/worker/bot/dbclient"
	"github.com/TicketsBot-cloud/worker/bot/logic"
)

const (
	preTicketProceedPrefix = "pretq_proceed_"
	preTicketRejectPrefix  = "pretq_reject_"
)

// panelFlowContext is satisfied by both ButtonContext and SelectMenuContext
type panelFlowContext interface {
	cmdregistry.InteractionContext
	Modal(res button.ResponseModal)
	Edit(data command.MessageResponse)
}

// startPanelFlow is the entrypoint after a panel button/select interaction has been validated. If the panel has
// pre-ticket questions configured, the first question is shown, and the ticket is not opened. Otherwise, the
// standard flow (form, or open directly) is continued.
func startPanelFlow(ctx panelFlowContext, panel database.Panel, outOfHoursTitle, outOfHoursWarning *string, outOfHoursColour *int) {
	questions, err := dbclient.Client.PanelPreTicketQuestions.GetByPanelId(ctx, panel.PanelId)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if len(questions) > 0 {
		res, err := buildPreTicketQuestionResponse(ctx, questions[0])
		if err != nil {
			ctx.HandleError(err)
			return
		}

		_, _ = ctx.ReplyWith(res)
		return
	}

	continuePanelFlow(ctx, panel, outOfHoursTitle, outOfHoursWarning, outOfHoursColour)
}

// continuePanelFlow opens the ticket, or shows the panel's form first if one is configured.
func continuePanelFlow(ctx panelFlowContext, panel database.Panel, outOfHoursTitle, outOfHoursWarning *string, outOfHoursColour *int) {
	if panel.FormId == nil {
		_, _ = logic.OpenTicket(ctx, ctx, &panel, panel.Title, nil, outOfHoursTitle, outOfHoursWarning, outOfHoursColour)
	} else {
		form, ok, err := dbclient.Client.Forms.Get(ctx, *panel.FormId)
		if err != nil {
			ctx.HandleError(err)
			return
		}

		if !ok {
			ctx.HandleError(errors.New("Form not found"))
			return
		}

		inputs, err := dbclient.Client.FormInput.GetInputs(ctx, form.Id)
		if err != nil {
			ctx.HandleError(err)
			return
		}

		inputOptions, err := dbclient.Client.FormInputOption.GetOptionsByForm(ctx, form.Id)
		if err != nil {
			ctx.HandleError(err)
			return
		}

		if len(inputs) == 0 { // Don't open a blank form
			_, _ = logic.OpenTicket(ctx, ctx, &panel, panel.Title, nil, outOfHoursTitle, outOfHoursWarning, outOfHoursColour)
		} else {
			modal := buildForm(panel, form, inputs, inputOptions)
			ctx.Modal(modal)
		}
	}
}

// buildPreTicketQuestionResponse builds the ephemeral message shown to the user for a pre-ticket question:
// the question embed, plus the proceed / reject buttons.
func buildPreTicketQuestionResponse(ctx cmdregistry.CommandContext, question database.PanelPreTicketQuestion) (command.MessageResponse, error) {
	embed, err := buildPreTicketEmbed(ctx, question.QuestionEmbedId)
	if err != nil {
		return command.MessageResponse{}, err
	}

	components := []component.Component{
		component.BuildActionRow(
			component.BuildButton(component.Button{
				Label:    question.ProceedButtonLabel,
				CustomId: fmt.Sprintf("%s%d", preTicketProceedPrefix, question.Id),
				Style:    component.ButtonStylePrimary,
			}),
			component.BuildButton(component.Button{
				Label:    question.RejectButtonLabel,
				CustomId: fmt.Sprintf("%s%d", preTicketRejectPrefix, question.Id),
				Style:    component.ButtonStyleSecondary,
			}),
		),
	}

	return command.NewEphemeralEmbedMessageResponseWithComponents(embed, components), nil
}

func buildPreTicketEmbed(ctx cmdregistry.CommandContext, embedId int) (*embed.Embed, error) {
	data, err := dbclient.Client.Embeds.GetEmbed(ctx, embedId)
	if err != nil {
		return nil, err
	}

	fields, err := dbclient.Client.EmbedFields.GetFieldsForEmbed(ctx, embedId)
	if err != nil {
		return nil, err
	}

	// No ticket exists at this point; a zero ticket ID skips placeholder substitutions, but the user ID is
	// still used for %avatar_url% image substitution
	pseudoTicket := database.Ticket{
		GuildId: ctx.GuildId(),
		UserId:  ctx.UserId(),
	}

	return logic.BuildCustomEmbed(ctx, ctx.Worker(), pseudoTicket, data, fields, false, nil), nil
}

// getQuestionAndPanel loads a pre-ticket question by ID and its panel, verifying it belongs to the interaction's guild.
func getQuestionAndPanel(ctx *context.ButtonContext, prefix string) (database.PanelPreTicketQuestion, database.Panel, bool) {
	questionId, err := strconv.Atoi(strings.TrimPrefix(ctx.InteractionData.CustomId, prefix))
	if err != nil {
		return database.PanelPreTicketQuestion{}, database.Panel{}, false
	}

	question, ok, err := dbclient.Client.PanelPreTicketQuestions.GetById(ctx, questionId)
	if err != nil {
		ctx.HandleError(err)
		return database.PanelPreTicketQuestion{}, database.Panel{}, false
	}

	if !ok { // Question was deleted since the message was sent
		return database.PanelPreTicketQuestion{}, database.Panel{}, false
	}

	panel, err := dbclient.Client.Panel.GetById(ctx, question.PanelId)
	if err != nil {
		ctx.HandleError(err)
		return database.PanelPreTicketQuestion{}, database.Panel{}, false
	}

	if panel.PanelId == 0 || panel.GuildId != ctx.GuildId() {
		return database.PanelPreTicketQuestion{}, database.Panel{}, false
	}

	return question, panel, true
}

type PreTicketProceedHandler struct{}

func (h *PreTicketProceedHandler) Matcher() matcher.Matcher {
	return matcher.NewFuncMatcher(func(customId string) bool {
		return strings.HasPrefix(customId, preTicketProceedPrefix)
	})
}

func (h *PreTicketProceedHandler) Properties() registry.Properties {
	return registry.Properties{
		Flags:   registry.SumFlags(registry.GuildAllowed, registry.CanEdit),
		Timeout: constants.TimeoutOpenTicket,
	}
}

func (h *PreTicketProceedHandler) Execute(ctx *context.ButtonContext) {
	question, panel, ok := getQuestionAndPanel(ctx, preTicketProceedPrefix)
	if !ok {
		return
	}

	// Re-validate panel access: support hours etc. may have changed since the question was shown
	canProceed, outOfHoursTitle, outOfHoursWarning, outOfHoursColour, err := logic.ValidatePanelAccess(ctx, panel)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if !canProceed {
		return
	}

	next, ok, err := dbclient.Client.PanelPreTicketQuestions.GetNext(ctx, question.PanelId, question.Position)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if ok { // Another question is configured: show it in place of the current one
		res, err := buildPreTicketQuestionResponse(ctx, next)
		if err != nil {
			ctx.HandleError(err)
			return
		}

		ctx.Edit(res)
		return
	}

	// No more questions: proceed to the form / open the ticket
	if panel.FormId == nil {
		// Remove the buttons from the question message before opening; OpenTicket's replies are sent as follow-ups
		embed, err := buildPreTicketEmbed(ctx, question.QuestionEmbedId)
		if err != nil {
			ctx.HandleError(err)
			return
		}

		ctx.Edit(command.NewEphemeralEmbedMessageResponse(embed))
	}

	continuePanelFlow(ctx, panel, outOfHoursTitle, outOfHoursWarning, outOfHoursColour)
}

type PreTicketRejectHandler struct{}

func (h *PreTicketRejectHandler) Matcher() matcher.Matcher {
	return matcher.NewFuncMatcher(func(customId string) bool {
		return strings.HasPrefix(customId, preTicketRejectPrefix)
	})
}

func (h *PreTicketRejectHandler) Properties() registry.Properties {
	return registry.Properties{
		Flags:   registry.SumFlags(registry.GuildAllowed, registry.CanEdit),
		Timeout: constants.TimeoutOpenTicket,
	}
}

func (h *PreTicketRejectHandler) Execute(ctx *context.ButtonContext) {
	question, _, ok := getQuestionAndPanel(ctx, preTicketRejectPrefix)
	if !ok {
		return
	}

	embed, err := buildPreTicketEmbed(ctx, question.RejectEmbedId)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	// Replace the question with the configured embed and remove the buttons; no ticket is opened
	ctx.Edit(command.NewEphemeralEmbedMessageResponse(embed))
}
