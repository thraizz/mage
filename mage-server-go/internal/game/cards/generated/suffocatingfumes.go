package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Suffocating Fumes", NewSuffocatingFumes)
}

// NewSuffocatingFumes creates a Suffocating Fumes
// {2}{B} - INSTANT
func NewSuffocatingFumes(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Suffocating Fumes")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
