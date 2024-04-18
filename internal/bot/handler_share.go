package bot

import (
	"context"
	"log"

	"github.com/dvordrova/feed_puppy/internal/database"
	"github.com/dvordrova/feed_puppy/internal/utils"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) OnShareDogForOwner(c tele.Context) error {
	ctx := context.TODO()

	u, _ := h.queries.GetUser(ctx, c.Sender().ID)
	hash := utils.GenerateBase64Hash()
	h.queries.NewInvite(ctx, database.NewInviteParams{
		DogID: u.CurrentDog.Int64,
		Hash:  hash,
		Type:  "owner",
	})
	dog, _ := h.queries.GetDog(ctx, u.CurrentDog.Int64)
	defer c.Bot().Delete(c.Message())
	return c.Send(h.layout.Text(c, "msg_share_dog_for_owner",
		struct {
			Hash    string
			DogName string
		}{
			Hash:    hash,
			DogName: dog.Name,
		},
		tele.Silent,
	))
}

func (h *Handler) OnShareDogForReader(c tele.Context) error {
	ctx := context.TODO()

	u, _ := h.queries.GetUser(ctx, c.Sender().ID)
	hash := utils.GenerateBase64Hash()
	h.queries.NewInvite(ctx, database.NewInviteParams{
		DogID: u.CurrentDog.Int64,
		Hash:  hash,
		Type:  "reader",
	})
	dog, _ := h.queries.GetDog(ctx, u.CurrentDog.Int64)
	defer c.Bot().Delete(c.Message())
	return c.Send(h.layout.Text(c, "msg_share_dog_for_reader",
		struct {
			Hash    string
			DogName string
		}{
			Hash:    hash,
			DogName: dog.Name,
		},
		tele.Silent,
	))
}

func (h *Handler) OnUnshareAllDogs(c tele.Context) error {
	log.Println("/unshare_all_dogs")
	defer c.Bot().Delete(c.Message())
	return c.Send(h.layout.Text(c, "msg_unshare_all_dogs"), tele.Silent)
}

func (h *Handler) OnUnsubscribe(c tele.Context) error {
	log.Println("/unsubscribe")
	defer c.Bot().Delete(c.Message())
	return c.Send(h.layout.Text(c, "msg_unsubscribe"), tele.Silent)
}
