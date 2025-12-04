package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bounty Of The Luxa", NewBountyOfTheLuxa)
}

// NewBountyOfTheLuxa creates a Bounty Of The Luxa
// {2}{G}{U} - ENCHANTMENT
func NewBountyOfTheLuxa(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bounty Of The Luxa")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
