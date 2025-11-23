package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Voice Of All", NewVoiceOfAll)
}

// NewVoiceOfAll creates a Voice Of All
// {2}{W}{W} - CREATURE
// Flying
func NewVoiceOfAll(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Voice Of All")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
