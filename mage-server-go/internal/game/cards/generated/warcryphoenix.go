package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Warcry Phoenix", NewWarcryPhoenix)
}

// NewWarcryPhoenix creates a Warcry Phoenix
// {3}{R} - CREATURE
// Flying, Haste
func NewWarcryPhoenix(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Warcry Phoenix")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHOENIX"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new ReturnToBattlefieldUnderOwnerControlSourceEffe...)
	// card.AddAbility(ability2)
	return card, nil
}
