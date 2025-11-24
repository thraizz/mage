package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Syndicate Heavy", NewSyndicateHeavy)
}

// NewSyndicateHeavy creates a Syndicate Heavy
// {2}{W/B}{W/B} - CREATURE
func NewSyndicateHeavy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Syndicate Heavy")
	card.ManaCost = "{2}{W/B}{W/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "ROGUE"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}