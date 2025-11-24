package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kroxa And Kunoros", NewKroxaAndKunoros)
}

// NewKroxaAndKunoros creates a Kroxa And Kunoros
// {3}{R}{W}{B} - CREATURE
// Vigilance, Lifelink
func NewKroxaAndKunoros(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kroxa And Kunoros")
	card.ManaCost = "{3}{R}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDER", "GIANT", "DOG"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	return card, nil
}