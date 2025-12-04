package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Command The Dreadhorde", NewCommandTheDreadhorde)
}

// NewCommandTheDreadhorde creates a Command The Dreadhorde
// {4}{B}{B} - SORCERY
func NewCommandTheDreadhorde(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Command The Dreadhorde")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
