package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nissa Who Shakes The World", NewNissaWhoShakesTheWorld)
}

// NewNissaWhoShakesTheWorld creates a Nissa Who Shakes The World
// {3}{G}{G} - PLANESWALKER
// Haste, Vigilance
func NewNissaWhoShakesTheWorld(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nissa Who Shakes The World")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NISSA", "ELEMENTAL"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "0"
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
