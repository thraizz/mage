package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abandon Attachments", NewAbandonAttachments)
}

// NewAbandonAttachments creates a Abandon Attachments
// {1}{U/R} - INSTANT
func NewAbandonAttachments(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abandon Attachments")
	card.ManaCost = "{1}{U/R}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(2), new Discard...)
	// card.AddAbility(ability0)
	return card, nil
}
