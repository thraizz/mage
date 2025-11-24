package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Atreus Impulsive Son", NewAtreusImpulsiveSon)
}

// NewAtreusImpulsiveSon creates a Atreus Impulsive Son
// {1}{U}{R} - CREATURE
// Reach
func NewAtreusImpulsiveSon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Atreus Impulsive Son")
	card.ManaCost = "{1}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOD", "ARCHER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - DiscardControllerEffect(1)
	//
	// Costs:
	//   - AddManaCost("{3}")
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
