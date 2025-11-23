package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Somberwald Stag", NewSomberwaldStag)
}

// NewSomberwaldStag creates a Somberwald Stag
// {3}{G}{G} - CREATURE
func NewSomberwaldStag(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Somberwald Stag")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELK"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
