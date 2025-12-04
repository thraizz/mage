package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Puresight Merrow", NewPuresightMerrow)
}

// NewPuresightMerrow creates a Puresight Merrow
// {W/U}{W/U} - CREATURE
func NewPuresightMerrow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Puresight Merrow")
	card.ManaCost = "{W/U}{W/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "WIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PuresightMerrowEffect()
	// card.AddAbility(ability0)
	return card, nil
}
