package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lamplight Phoenix", NewLamplightPhoenix)
}

// NewLamplightPhoenix creates a Lamplight Phoenix
// {1}{R}{R} - CREATURE
// Flying
func NewLamplightPhoenix(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lamplight Phoenix")
	card.ManaCost = "{1}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHOENIX"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new ReturnToBattlefieldUnderOwner...)
	// card.AddAbility(ability1)
	return card, nil
}
