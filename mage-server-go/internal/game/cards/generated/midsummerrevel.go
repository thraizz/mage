package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Midsummer Revel", NewMidsummerRevel)
}

// NewMidsummerRevel creates a Midsummer Revel
// {3}{G}{G} - ENCHANTMENT
func NewMidsummerRevel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Midsummer Revel")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("BeastToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewCreateTokenEffect(token0_0, new CountersSourceCount(CounterType.VERSE))).
		Build()
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeVerse.CreateInstance(1), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}