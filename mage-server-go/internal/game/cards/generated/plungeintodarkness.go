package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Plunge Into Darkness", NewPlungeIntoDarkness)
}

// NewPlungeIntoDarkness creates a Plunge Into Darkness
// {1}{B} - INSTANT
func NewPlungeIntoDarkness(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Plunge Into Darkness")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
