package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Deity Of Scars", NewDeityOfScars)
}

// NewDeityOfScars creates a Deity Of Scars
// {B/G}{B/G}{B/G}{B/G}{B/G} - CREATURE
// Trample
func NewDeityOfScars(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Deity Of Scars")
	card.ManaCost = "{B/G}{B/G}{B/G}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "AVATAR"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - RegenerateSourceEffect()
	// card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeM1M1.CreateInstance(2))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
