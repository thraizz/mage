package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bringer Of The Black Dawn", NewBringerOfTheBlackDawn)
}

// NewBringerOfTheBlackDawn creates a Bringer Of The Black Dawn
// {7}{B}{B} - CREATURE
// Trample
func NewBringerOfTheBlackDawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bringer Of The Black Dawn")
	card.ManaCost = "{7}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BRINGER"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfUpkeepTriggeredAbility
	//   - Effect: DoIfCostPaid(new SearchLibraryPutOnLibraryEffect(new TargetCard...)
	// card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new SearchLibraryPutOnLibraryEffect(new TargetCard...)
	// card.AddAbility(ability2)
	return card, nil
}
