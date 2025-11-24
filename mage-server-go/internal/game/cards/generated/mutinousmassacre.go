package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mutinous Massacre", NewMutinousMassacre)
}

// NewMutinousMassacre creates a Mutinous Massacre
// {3}{B}{B}{R}{R} - SORCERY
func NewMutinousMassacre(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mutinous Massacre")
	card.ManaCost = "{3}{B}{B}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}