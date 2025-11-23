package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unearthly Blizzard", NewUnearthlyBlizzard)
}

// NewUnearthlyBlizzard creates a Unearthly Blizzard
// {2}{R} - SORCERY
func NewUnearthlyBlizzard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unearthly Blizzard")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
