package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Formless Genesis", NewFormlessGenesis)
}

// NewFormlessGenesis creates a Formless Genesis
// {2}{G} - KINDRED SORCERY
func NewFormlessGenesis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Formless Genesis")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}