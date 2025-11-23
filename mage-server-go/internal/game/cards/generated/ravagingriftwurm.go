package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Ravaging Riftwurm", NewRavagingRiftwurm)
}

// NewRavagingRiftwurm creates a Ravaging Riftwurm
// {1}{G}{G} - CREATURE
func NewRavagingRiftwurm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ravaging Riftwurm")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WURM"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeTime.CreateInstance(3))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
