package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shell Of The Last Kappa", NewShellOfTheLastKappa)
}

// NewShellOfTheLastKappa creates a Shell Of The Last Kappa
// {3} - ARTIFACT
func NewShellOfTheLastKappa(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shell Of The Last Kappa")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
