package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaseto Orochi Archmage", NewKasetoOrochiArchmage)
}

// NewKasetoOrochiArchmage creates a Kaseto Orochi Archmage
// {1}{G}{U} - CREATURE
func NewKasetoOrochiArchmage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaseto Orochi Archmage")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - KasetoEffect()
	// card.AddAbility(ability0)
	return card, nil
}
