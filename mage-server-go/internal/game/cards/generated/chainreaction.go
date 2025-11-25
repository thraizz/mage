package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chain Reaction", NewChainReaction)
}

// NewChainReaction creates a Chain Reaction
// {2}{R}{R} - SORCERY
func NewChainReaction(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chain Reaction")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(new PermanentsOnBattlefieldCount(new FilterCreatur...)
	// card.AddAbility(ability0)
	return card, nil
}
