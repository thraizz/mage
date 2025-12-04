package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Greasewrench Goblin", NewGreasewrenchGoblin)
}

// NewGreasewrenchGoblin creates a Greasewrench Goblin
// {R} - CREATURE
func NewGreasewrenchGoblin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Greasewrench Goblin")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "ARTIFICER"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ExhaustAbility
	//   - Effect: DiscardAndDrawThatManyEffect()
	// card.AddAbility(ability0)
	return card, nil
}
