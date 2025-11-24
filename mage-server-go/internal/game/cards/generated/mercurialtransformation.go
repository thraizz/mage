package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mercurial Transformation", NewMercurialTransformation)
}

// NewMercurialTransformation creates a Mercurial Transformation
// {1}{U} - SORCERY
func NewMercurialTransformation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mercurial Transformation")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}