package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ideas Unbound", NewIdeasUnbound)
}

// NewIdeasUnbound creates a Ideas Unbound
// {U}{U} - SORCERY
func NewIdeasUnbound(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ideas Unbound")
	card.ManaCost = "{U}{U}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(3)
	// card.AddAbility(ability0)
	return card, nil
}
