package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crystal Shard", NewCrystalShard)
}

// NewCrystalShard creates a Crystal Shard
// {3} - ARTIFACT
func NewCrystalShard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crystal Shard")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - CrystalShardEffect()
	//
	// Costs:
	//   - AddTapCost()
	//   - AddManaCost("{3}")
	// card.AddAbility(ability0)
	return card, nil
}
