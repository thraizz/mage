package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Frilled Cave Wurm", NewFrilledCaveWurm)
}

// NewFrilledCaveWurm creates a Frilled Cave Wurm
// {3}{U} - CREATURE
func NewFrilledCaveWurm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Frilled Cave Wurm")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SALAMANDER", "WURM"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(2, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
