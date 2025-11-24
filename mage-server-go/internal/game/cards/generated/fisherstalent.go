package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fishers Talent", NewFishersTalent)
}

// NewFishersTalent creates a Fishers Talent
// {2}{G}{U} - ENCHANTMENT
func NewFishersTalent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fishers Talent")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CLASS"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
