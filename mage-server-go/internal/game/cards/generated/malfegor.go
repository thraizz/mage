package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Malfegor", NewMalfegor)
}

// NewMalfegor creates a Malfegor
// {2}{B}{B}{R}{R} - CREATURE
// Flying
func NewMalfegor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Malfegor")
	card.ManaCost = "{2}{B}{B}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON", "DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
