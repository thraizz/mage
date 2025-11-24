package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bonfire Of The Damned", NewBonfireOfTheDamned)
}

// NewBonfireOfTheDamned creates a Bonfire Of The Damned
// {X}{X}{R} - SORCERY
func NewBonfireOfTheDamned(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bonfire Of The Damned")
	card.ManaCost = "{X}{X}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
