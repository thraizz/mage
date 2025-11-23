package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Impromptu Raid", NewImpromptuRaid)
}

// NewImpromptuRaid creates a Impromptu Raid
// {3}{R/G} - ENCHANTMENT
func NewImpromptuRaid(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Impromptu Raid")
	card.ManaCost = "{3}{R/G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("", source.getControllerId())
	// card.AddAbility(ability0)
	return card, nil
}
