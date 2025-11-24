package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dromar The Banisher", NewDromarTheBanisher)
}

// NewDromarTheBanisher creates a Dromar The Banisher
// {3}{W}{U}{B} - CREATURE
// Flying
func NewDromarTheBanisher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dromar The Banisher")
	card.ManaCost = "{3}{W}{U}{B}"
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
	//   - DoIfCostPaid(new DromarTheBanisherEffect(), new ManaCostsImpl<>...)
	// card.AddAbility(ability1)
	return card, nil
}