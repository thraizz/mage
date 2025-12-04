package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arcane Endeavor", NewArcaneEndeavor)
}

// NewArcaneEndeavor creates a Arcane Endeavor
// {5}{U}{U} - SORCERY
func NewArcaneEndeavor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arcane Endeavor")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
