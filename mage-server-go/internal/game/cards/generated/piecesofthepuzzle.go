package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pieces Of The Puzzle", NewPiecesOfThePuzzle)
}

// NewPiecesOfThePuzzle creates a Pieces Of The Puzzle
// {2}{U} - SORCERY
func NewPiecesOfThePuzzle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pieces Of The Puzzle")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(5, 2, filter, PutCards.HAND, PutCards.GRAVEYARD, f...)
	// card.AddAbility(ability0)
	return card, nil
}
