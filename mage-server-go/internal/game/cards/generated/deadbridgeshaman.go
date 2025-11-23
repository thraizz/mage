package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Deadbridge Shaman", NewDeadbridgeShaman)
}

// NewDeadbridgeShaman creates a Deadbridge Shaman
// {2}{B} - CREATURE
func NewDeadbridgeShaman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Deadbridge Shaman")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "SHAMAN"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	// card.AddAbility(ability0)
	return card, nil
}
