package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Innkeepers Talent", NewInnkeepersTalent)
}

// NewInnkeepersTalent creates a Innkeepers Talent
// {1}{G} - ENCHANTMENT
func NewInnkeepersTalent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Innkeepers Talent")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CLASS"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect(new WardAbility(, false), filter)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}