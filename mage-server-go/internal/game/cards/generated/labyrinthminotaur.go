package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Labyrinth Minotaur", NewLabyrinthMinotaur)
}

// NewLabyrinthMinotaur creates a Labyrinth Minotaur
// {3}{U} - CREATURE
func NewLabyrinthMinotaur(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Labyrinth Minotaur")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MINOTAUR"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
