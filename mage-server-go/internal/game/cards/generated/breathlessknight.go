package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Breathless Knight", NewBreathlessKnight)
}

// NewBreathlessKnight creates a Breathless Knight
// {1}{W}{B} - CREATURE
// Flying, Lifelink
func NewBreathlessKnight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Breathless Knight")
	card.ManaCost = "{1}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "KNIGHT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	return card, nil
}
