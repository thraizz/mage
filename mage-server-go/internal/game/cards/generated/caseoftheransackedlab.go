package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Case Of The Ransacked Lab", NewCaseOfTheRansackedLab)
}

// NewCaseOfTheRansackedLab creates a Case Of The Ransacked Lab
// {2}{U} - ENCHANTMENT
func NewCaseOfTheRansackedLab(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Case Of The Ransacked Lab")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CASE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
