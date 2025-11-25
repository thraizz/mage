package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Laurine The Diversion", NewLaurineTheDiversion)
}

// NewLaurineTheDiversion creates a Laurine The Diversion
// {2}{R} - CREATURE
// FirstStrike
func NewLaurineTheDiversion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Laurine The Diversion")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - GoadTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	// card.AddAbility(ability1)
	return card, nil
}
