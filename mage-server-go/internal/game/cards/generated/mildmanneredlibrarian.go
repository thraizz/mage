package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mild Mannered Librarian", NewMildManneredLibrarian)
}

// NewMildManneredLibrarian creates a Mild Mannered Librarian
// {G} - CREATURE
func NewMildManneredLibrarian(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mild Mannered Librarian")
	card.ManaCost = "{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
