package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crookshank Kobolds", NewCrookshankKobolds)
}

// NewCrookshankKobolds creates a Crookshank Kobolds
// {0} - CREATURE
func NewCrookshankKobolds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crookshank Kobolds")
	card.ManaCost = "{0}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOBOLD"}
	card.Power = "0"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}