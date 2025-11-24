package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fell Beasts Shriek", NewFellBeastsShriek)
}

// NewFellBeastsShriek creates a Fell Beasts Shriek
// {U}{R} - SORCERY
func NewFellBeastsShriek(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fell Beasts Shriek")
	card.ManaCost = "{U}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}