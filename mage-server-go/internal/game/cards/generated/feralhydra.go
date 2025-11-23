package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Feral Hydra", NewFeralHydra)
}

// NewFeralHydra creates a Feral Hydra
// {X}{G} - CREATURE
func NewFeralHydra(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Feral Hydra")
	card.ManaCost = "{X}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HYDRA", "BEAST"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
