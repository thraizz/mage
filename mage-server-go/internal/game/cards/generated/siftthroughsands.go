package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sift Through Sands", NewSiftThroughSands)
}

// NewSiftThroughSands creates a Sift Through Sands
// {1}{U}{U} - INSTANT
func NewSiftThroughSands(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sift Through Sands")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(1)
	//   - DiscardControllerEffect(1)
	// card.AddAbility(ability0)
	return card, nil
}
