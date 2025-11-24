package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ozolith The Shattered Spire", NewOzolithTheShatteredSpire)
}

// NewOzolithTheShatteredSpire creates a Ozolith The Shattered Spire
// {1}{G} - ARTIFACT
func NewOzolithTheShatteredSpire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ozolith The Shattered Spire")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}