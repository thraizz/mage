package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Firesong And Sunspeaker", NewFiresongAndSunspeaker)
}

// NewFiresongAndSunspeaker creates a Firesong And Sunspeaker
// {4}{R}{W} - CREATURE
// Lifelink
func NewFiresongAndSunspeaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Firesong And Sunspeaker")
	card.ManaCost = "{4}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MINOTAUR", "CLERIC"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	return card, nil
}
