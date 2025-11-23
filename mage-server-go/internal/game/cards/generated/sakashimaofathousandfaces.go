package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sakashima Of A Thousand Faces", NewSakashimaOfAThousandFaces)
}

// NewSakashimaOfAThousandFaces creates a Sakashima Of A Thousand Faces
// {3}{U} - CREATURE
func NewSakashimaOfAThousandFaces(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sakashima Of A Thousand Faces")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(                 StaticFilters.FILTER_CONTROLLED_A...)
	// card.AddAbility(ability0)
	return card, nil
}
