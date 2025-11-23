package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Host Of The Hereafter", NewHostOfTheHereafter)
}

// NewHostOfTheHereafter creates a Host Of The Hereafter
// {2}{B}{G} - CREATURE
func NewHostOfTheHereafter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Host Of The Hereafter")
	card.ManaCost = "{2}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "WARLOCK"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(2))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
