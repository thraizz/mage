package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ulalek Fused Atrocity", NewUlalekFusedAtrocity)
}

// NewUlalekFusedAtrocity creates a Ulalek Fused Atrocity
// {C/W}{C/U}{C/B}{C/R}{C/G} - CREATURE
func NewUlalekFusedAtrocity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ulalek Fused Atrocity")
	card.ManaCost = "{C/W}{C/U}{C/B}{C/R}{C/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new UlalekFusedAtrocityEffect(), new ColorlessMana...)
	// card.AddAbility(ability0)
	return card, nil
}
