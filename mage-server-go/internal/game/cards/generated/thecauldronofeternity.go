package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Cauldron Of Eternity", NewTheCauldronOfEternity)
}

// NewTheCauldronOfEternity creates a The Cauldron Of Eternity
// {10}{B}{B} - ARTIFACT
func NewTheCauldronOfEternity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Cauldron Of Eternity")
	card.ManaCost = "{10}{B}{B}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
