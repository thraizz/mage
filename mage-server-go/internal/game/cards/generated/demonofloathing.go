package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Demon Of Loathing", NewDemonOfLoathing)
}

// NewDemonOfLoathing creates a Demon Of Loathing
// {5}{B}{B} - CREATURE
// Flying, Trample
func NewDemonOfLoathing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Demon Of Loathing")
	card.ManaCost = "{5}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	return card, nil
}
