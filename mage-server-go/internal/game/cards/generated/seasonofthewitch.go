package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Season Of The Witch", NewSeasonOfTheWitch)
}

// NewSeasonOfTheWitch creates a Season Of The Witch
// {B}{B}{B} - ENCHANTMENT
func NewSeasonOfTheWitch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Season Of The Witch")
	card.ManaCost = "{B}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: SeasonOfTheWitchEffect()
	// card.AddAbility(ability0)
	return card, nil
}
