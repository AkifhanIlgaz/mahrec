package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"mahrec/internal/middleware"
	"mahrec/internal/store"
	"mahrec/internal/web/templates/components"
	"mahrec/internal/web/templates/pages"
)

type GoalHandler struct {
	goals *store.GoalStore
}

func NewGoalHandler(goals *store.GoalStore) *GoalHandler {
	return &GoalHandler{goals: goals}
}

func (h *GoalHandler) Index(c fiber.Ctx) error {
	clientID := middleware.FromCtx(c)

	goals, err := h.goals.List(clientID)
	if err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return pages.Index(goals).Render(c.Context(), c.Response().BodyWriter())
}

func (h *GoalHandler) Create(c fiber.Ctx) error {
	clientID := middleware.FromCtx(c)

	title := c.FormValue("title")
	if title == "" {
		return renderFormError(c, "Listeden okunacak bir şey seçmelisin.")
	}

	rawTarget := c.FormValue("target_count")
	target, err := strconv.Atoi(rawTarget)
	if err != nil || target <= 0 {
		return renderFormError(c, "Adet, 1 veya daha büyük bir sayı olmalı.")
	}

	description := c.FormValue("description")

	goal, err := h.goals.Create(clientID, title, description, target)
	if err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	c.Set("HX-Trigger", "goal-created")
	if err := components.GoalListItem(goal).Render(c.Context(), c.Response().BodyWriter()); err != nil {
		return err
	}
	if err := components.FormError("").Render(c.Context(), c.Response().BodyWriter()); err != nil {
		return err
	}
	// Hide the static "no goals yet" paragraph, which only server-renders
	// based on the list at page-load time and wouldn't otherwise disappear
	// after the first goal is added via htmx.
	_, err = c.Response().BodyWriter().Write([]byte(`<p id="empty-state" hidden hx-swap-oob="true"></p>`))
	return err
}

// renderFormError writes the new-goal form's error banner as an
// out-of-band swap so the user sees why their submission was rejected.
func renderFormError(c fiber.Ctx, message string) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return components.FormError(message).Render(c.Context(), c.Response().BodyWriter())
}

func (h *GoalHandler) Increment(c fiber.Ctx) error {
	clientID := middleware.FromCtx(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid goal id")
	}

	amount, err := strconv.Atoi(c.FormValue("amount"))
	if err != nil || amount <= 0 {
		amount = 1
	}

	goal, err := h.goals.Increment(clientID, id, amount)
	if err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return components.GoalItem(goal).Render(c.Context(), c.Response().BodyWriter())
}

func (h *GoalHandler) Reset(c fiber.Ctx) error {
	clientID := middleware.FromCtx(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid goal id")
	}

	goal, err := h.goals.Reset(clientID, id)
	if err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return components.GoalItem(goal).Render(c.Context(), c.Response().BodyWriter())
}

func (h *GoalHandler) Delete(c fiber.Ctx) error {
	clientID := middleware.FromCtx(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid goal id")
	}

	if err := h.goals.Delete(clientID, id); err != nil {
		return err
	}

	c.Status(fiber.StatusOK)
	return c.SendString("")
}
