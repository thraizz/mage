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
	cards.Register("Phyrexian Swarmlord", NewPhyrexianSwarmlord)
}

// NewPhyrexianSwarmlord creates a Phyrexian Swarmlord
// {4}{G}{G} - CREATURE
func NewPhyrexianSwarmlord(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexian Swarmlord")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "INSECT", "HORROR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("InsectInfectToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, OpponentsPoisonCountersCount.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
