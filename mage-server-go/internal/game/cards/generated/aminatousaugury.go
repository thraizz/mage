package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aminatous Augury", NewAminatousAugury)
}

// NewAminatousAugury creates a Aminatous Augury
// {6}{U}{U} - SORCERY
func NewAminatousAugury(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aminatous Augury")
	card.ManaCost = "{6}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
