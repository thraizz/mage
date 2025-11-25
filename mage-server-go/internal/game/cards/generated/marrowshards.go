package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Marrow Shards", NewMarrowShards)
}

// NewMarrowShards creates a Marrow Shards
// {W/P} - INSTANT
func NewMarrowShards(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Marrow Shards")
	card.ManaCost = "{W/P}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(1, new FilterAttackingCreature())
	// card.AddAbility(ability0)
	return card, nil
}
