package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vedalken Entrancer", NewVedalkenEntrancer)
}

// NewVedalkenEntrancer creates a Vedalken Entrancer
// {3}{U} - CREATURE
func NewVedalkenEntrancer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vedalken Entrancer")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VEDALKEN", "WIZARD"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewMillCardsTargetEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
