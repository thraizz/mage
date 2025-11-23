package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eyetwitch", NewEyetwitch)
}

// NewEyetwitch creates a Eyetwitch
// {B} - CREATURE
// Flying
func NewEyetwitch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eyetwitch")
	card.ManaCost = "{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"EYE", "BAT"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
