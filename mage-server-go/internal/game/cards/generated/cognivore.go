package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cognivore", NewCognivore)
}

// NewCognivore creates a Cognivore
// {6}{U}{U} - CREATURE
// Flying
func NewCognivore(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cognivore")
	card.ManaCost = "{6}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LHURGOYF"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
