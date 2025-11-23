package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Mishras Self Replicator", NewMishrasSelfReplicator)
}

// NewMishrasSelfReplicator creates a Mishras Self Replicator
// {5} - ARTIFACT CREATURE
func NewMishrasSelfReplicator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mishras Self Replicator")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ASSEMBLY_WORKER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateTokenCopySourceEffect()...)
	// card.AddAbility(ability0)
	return card, nil
}
