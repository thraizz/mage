package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Altar Of The Wretched", NewAltarOfTheWretched)
}

// NewAltarOfTheWretched creates a Altar Of The Wretched
// {2}{B} - ARTIFACT
func NewAltarOfTheWretched(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Altar Of The Wretched")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}