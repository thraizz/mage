package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mirkwood Trapper", NewMirkwoodTrapper)
}

// NewMirkwoodTrapper creates a Mirkwood Trapper
// {1}{G}{U} - CREATURE
func NewMirkwoodTrapper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mirkwood Trapper")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "SCOUT"}
	card.Power = "1"
	card.Toughness = "4"
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
