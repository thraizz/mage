package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Seraph Of The Sword", NewSeraphOfTheSword)
}

// NewSeraphOfTheSword creates a Seraph Of The Sword
// {3}{W} - CREATURE
// Flying
func NewSeraphOfTheSword(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Seraph Of The Sword")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
