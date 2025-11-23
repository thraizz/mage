package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Killer Instinct", NewKillerInstinct)
}

// NewKillerInstinct creates a Killer Instinct
// {4}{R}{G} - ENCHANTMENT
func NewKillerInstinct(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Killer Instinct")
	card.ManaCost = "{4}{R}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("Sacrifice it at the beginning of the next end ste...)
	// card.AddAbility(ability0)
	return card, nil
}
