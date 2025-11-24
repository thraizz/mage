package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sigarda Champion Of Light", NewSigardaChampionOfLight)
}

// NewSigardaChampionOfLight creates a Sigarda Champion Of Light
// {1}{G}{W}{W} - CREATURE
// Flying, Trample
func NewSigardaChampionOfLight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sigarda Champion Of Light")
	card.ManaCost = "{1}{G}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: BoostControlledEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 5, 1, filter2, PutCards.HAND, Put...)
	// card.AddAbility(ability3)
	return card, nil
}
