package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Go Shintai Of Shared Purpose", NewGoShintaiOfSharedPurpose)
}

// NewGoShintaiOfSharedPurpose creates a Go Shintai Of Shared Purpose
// {3}{W} - ENCHANTMENT CREATURE
// Vigilance
func NewGoShintaiOfSharedPurpose(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Go Shintai Of Shared Purpose")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"SHRINE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new CreateTokenEffect(new...)
	// card.AddAbility(ability1)
	return card, nil
}
