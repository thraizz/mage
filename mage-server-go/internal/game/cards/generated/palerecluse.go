package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pale Recluse", NewPaleRecluse)
}

// NewPaleRecluse creates a Pale Recluse
// {4}{G}{W} - CREATURE
// Reach
func NewPaleRecluse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pale Recluse")
	card.ManaCost = "{4}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIDER"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	return card, nil
}
