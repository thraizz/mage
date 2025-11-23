package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sky Eel School", NewSkyEelSchool)
}

// NewSkyEelSchool creates a Sky Eel School
// {3}{U}{U} - CREATURE
// Flying
func NewSkyEelSchool(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sky Eel School")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FISH"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
