package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Catacomb Slug", NewCatacombSlug)
}

// NewCatacombSlug creates a Catacomb Slug
// {4}{B} - CREATURE
func NewCatacombSlug(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Catacomb Slug")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLUG"}
	card.Power = "2"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
