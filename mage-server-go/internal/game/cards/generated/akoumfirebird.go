package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Akoum Firebird", NewAkoumFirebird)
}

// NewAkoumFirebird creates a Akoum Firebird
// {2}{R}{R} - CREATURE
// Flying, Haste
func NewAkoumFirebird(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Akoum Firebird")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHOENIX"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new ReturnSourceFromGraveyardToBa...)
	// card.AddAbility(ability2)
	return card, nil
}
