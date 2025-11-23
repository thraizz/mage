package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Esoteric Duplicator", NewEsotericDuplicator)
}

// NewEsotericDuplicator creates a Esoteric Duplicator
// {2}{U} - ARTIFACT
func NewEsotericDuplicator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Esoteric Duplicator")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"CLUE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new EsotericDuplicatorEffect(), new GenericManaCos...)
	// card.AddAbility(ability0)
	return card, nil
}
