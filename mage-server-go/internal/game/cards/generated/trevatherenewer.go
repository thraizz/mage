package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Treva The Renewer", NewTrevaTheRenewer)
}

// NewTrevaTheRenewer creates a Treva The Renewer
// {3}{G}{W}{U} - CREATURE
// Flying
func NewTrevaTheRenewer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Treva The Renewer")
	card.ManaCost = "{3}{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new TrevaTheRenewerEffect(), new ManaCostsImpl<>("...)
	// card.AddAbility(ability1)
	return card, nil
}