package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Oko The Trickster", NewOkoTheTrickster)
}

// NewOkoTheTrickster creates a Oko The Trickster
// {4}{G}{U} - PLANESWALKER
func NewOkoTheTrickster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Oko The Trickster")
	card.ManaCost = "{4}{G}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"OKO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}