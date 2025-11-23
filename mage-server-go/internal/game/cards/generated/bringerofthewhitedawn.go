package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bringer Of The White Dawn", NewBringerOfTheWhiteDawn)
}

// NewBringerOfTheWhiteDawn creates a Bringer Of The White Dawn
// {7}{W}{W} - CREATURE
// Trample
func NewBringerOfTheWhiteDawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bringer Of The White Dawn")
	card.ManaCost = "{7}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BRINGER"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability1)
	return card, nil
}
