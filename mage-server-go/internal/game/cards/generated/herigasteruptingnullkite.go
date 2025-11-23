package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Herigast Erupting Nullkite", NewHerigastEruptingNullkite)
}

// NewHerigastEruptingNullkite creates a Herigast Erupting Nullkite
// {9} - CREATURE
// Flying
func NewHerigastEruptingNullkite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Herigast Erupting Nullkite")
	card.ManaCost = "{9}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI", "DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new DrawCardSourceControllerEffec...)
	// card.AddAbility(ability1)
	return card, nil
}
