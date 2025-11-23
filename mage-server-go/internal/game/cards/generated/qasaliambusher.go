package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Qasali Ambusher", NewQasaliAmbusher)
}

// NewQasaliAmbusher creates a Qasali Ambusher
// {1}{G}{W} - CREATURE
// Reach
func NewQasaliAmbusher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Qasali Ambusher")
	card.ManaCost = "{1}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "WARRIOR"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	return card, nil
}
