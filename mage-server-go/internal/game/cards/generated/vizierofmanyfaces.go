package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vizier Of Many Faces", NewVizierOfManyFaces)
}

// NewVizierOfManyFaces creates a Vizier Of Many Faces
// {2}{U}{U} - CREATURE
func NewVizierOfManyFaces(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vizier Of Many Faces")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER", "CLERIC"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(new VizierOfManyFacesCopyApplier())
	// card.AddAbility(ability0)
	return card, nil
}