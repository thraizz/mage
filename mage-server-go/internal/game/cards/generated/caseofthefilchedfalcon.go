package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Case Of The Filched Falcon", NewCaseOfTheFilchedFalcon)
}

// NewCaseOfTheFilchedFalcon creates a Case Of The Filched Falcon
// {U} - ENCHANTMENT
// Flying
func NewCaseOfTheFilchedFalcon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Case Of The Filched Falcon")
	card.ManaCost = "{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CASE", "BIRD"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
