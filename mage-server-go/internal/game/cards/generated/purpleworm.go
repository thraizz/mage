package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Purple Worm", NewPurpleWorm)
}

// NewPurpleWorm creates a Purple Worm
// {5}{G}{G} - CREATURE
func NewPurpleWorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Purple Worm")
	card.ManaCost = "{5}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WORM"}
	card.Power = "8"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
