package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Juvenile Gloomwidow", NewJuvenileGloomwidow)
}

// NewJuvenileGloomwidow creates a Juvenile Gloomwidow
// {G}{G} - CREATURE
// Reach
func NewJuvenileGloomwidow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Juvenile Gloomwidow")
	card.ManaCost = "{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIDER"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	return card, nil
}
