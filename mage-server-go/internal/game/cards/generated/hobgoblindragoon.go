package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hobgoblin Dragoon", NewHobgoblinDragoon)
}

// NewHobgoblinDragoon creates a Hobgoblin Dragoon
// {2}{R/W} - CREATURE
// Flying, FirstStrike
func NewHobgoblinDragoon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hobgoblin Dragoon")
	card.ManaCost = "{2}{R/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "KNIGHT"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability1)
	return card, nil
}
