package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tergrids Shadow", NewTergridsShadow)
}

// NewTergridsShadow creates a Tergrids Shadow
// {3}{B}{B} - INSTANT
func NewTergridsShadow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tergrids Shadow")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(2, filter)
	// card.AddAbility(ability0)
	return card, nil
}
