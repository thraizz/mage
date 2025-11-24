package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ureni Of The Unwritten", NewUreniOfTheUnwritten)
}

// NewUreniOfTheUnwritten creates a Ureni Of The Unwritten
// {4}{G}{U}{R} - CREATURE
// Flying, Trample
func NewUreniOfTheUnwritten(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ureni Of The Unwritten")
	card.ManaCost = "{4}{G}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(8, 1, filter, PutCards.BATTLEFIELD, PutCards.BOTTO...)
	// card.AddAbility(ability2)
	return card, nil
}