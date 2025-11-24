package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Marwyns Kindred", NewMarwynsKindred)
}

// NewMarwynsKindred creates a Marwyns Kindred
// {X}{2}{G}{G} - SORCERY
func NewMarwynsKindred(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Marwyns Kindred")
	card.ManaCost = "{X}{2}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}