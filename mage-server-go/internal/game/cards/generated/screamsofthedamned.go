package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Screams Of The Damned", NewScreamsOfTheDamned)
}

// NewScreamsOfTheDamned creates a Screams Of The Damned
// {3}{B}{B} - ENCHANTMENT
func NewScreamsOfTheDamned(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Screams Of The Damned")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - DamageEverythingEffect()
	// card.AddAbility(ability0)
	return card, nil
}
