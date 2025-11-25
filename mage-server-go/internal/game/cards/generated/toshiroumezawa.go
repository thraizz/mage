package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Toshiro Umezawa", NewToshiroUmezawa)
}

// NewToshiroUmezawa creates a Toshiro Umezawa
// {1}{B}{B} - CREATURE
func NewToshiroUmezawa(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Toshiro Umezawa")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SAMURAI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: DiesCreatureTriggeredAbility
	//   - Effect: MayCastTargetCardEffect()
	// card.AddAbility(ability0)
	return card, nil
}
