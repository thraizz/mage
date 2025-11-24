package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vronos Masked Inquisitor", NewVronosMaskedInquisitor)
}

// NewVronosMaskedInquisitor creates a Vronos Masked Inquisitor
// {3}{U}{U} - PLANESWALKER
func NewVronosMaskedInquisitor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vronos Masked Inquisitor")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"VRONOS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - PhaseOutTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
