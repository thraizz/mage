package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Practical Research", NewPracticalResearch)
}

// NewPracticalResearch creates a Practical Research
// {3}{U}{R} - INSTANT
func NewPracticalResearch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Practical Research")
	card.ManaCost = "{3}{U}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(2)
	// card.AddAbility(ability0)
	return card, nil
}
