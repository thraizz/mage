package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pearl Shard", NewPearlShard)
}

// NewPearlShard creates a Pearl Shard
// {3} - ARTIFACT
func NewPearlShard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pearl Shard")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PreventDamageToTargetEffect()
	//
	// Costs:
	//   - AddTapCost()
	//   - AddManaCost("{3}")
	// card.AddAbility(ability0)
	return card, nil
}
