package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kokusho The Evening Star", NewKokushoTheEveningStar)
}

// NewKokushoTheEveningStar creates a Kokusho The Evening Star
// {4}{B}{B} - CREATURE
// Flying
func NewKokushoTheEveningStar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kokusho The Evening Star")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON", "SPIRIT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
