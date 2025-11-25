package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Quintorius Loremaster", NewQuintoriusLoremaster)
}

// NewQuintoriusLoremaster creates a Quintorius Loremaster
// {3}{R}{W} - CREATURE
// Vigilance
func NewQuintoriusLoremaster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quintorius Loremaster")
	card.ManaCost = "{3}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEPHANT", "CLERIC"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: ExileTargetForSourceEffect()
	// card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - QuintoriusLoremasterEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability2)
	return card, nil
}
