package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Treetop Sentinel", NewTreetopSentinel)
}

// NewTreetopSentinel creates a Treetop Sentinel
// {2}{U}{U} - CREATURE
// Flying
func NewTreetopSentinel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Treetop Sentinel")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD", "SOLDIER"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
