package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Supreme Inquisitor", NewSupremeInquisitor)
}

// NewSupremeInquisitor creates a Supreme Inquisitor
// {3}{U}{U} - CREATURE
func NewSupremeInquisitor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Supreme Inquisitor")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - SearchLibraryAndExileTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
