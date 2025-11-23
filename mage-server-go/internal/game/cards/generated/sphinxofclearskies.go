package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sphinx Of Clear Skies", NewSphinxOfClearSkies)
}

// NewSphinxOfClearSkies creates a Sphinx Of Clear Skies
// {3}{U}{U} - CREATURE
// Flying
func NewSphinxOfClearSkies(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sphinx Of Clear Skies")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPHINX"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
