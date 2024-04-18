package bot

import (
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) OnActionCommand(c tele.Context) error {
	defer c.Bot().Delete(c.Message())

	return c.Send(
		h.layout.Text(c, "msg_choose_action"),
		h.layout.Markup(c, "action"),
		tele.Silent,
	)
}

func (h *Handler) OnNewDog(c tele.Context) error {
	defer c.Bot().Delete(c.Message())
	return c.Send(h.layout.Text(c, "msg_new_dog"), tele.Silent)
}

func (h *Handler) OnChangeCurDog(c tele.Context) error {
	defer c.Bot().Delete(c.Message())
	return c.Send(h.layout.Text(c, "msg_change_cur_dog"), tele.Silent)
}
