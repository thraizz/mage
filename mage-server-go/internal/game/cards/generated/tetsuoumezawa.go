package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tetsuo Umezawa", NewTetsuoUmezawa)
}

// NewTetsuoUmezawa creates a Tetsuo Umezawa
// {U}{B}{R} - CREATURE
func NewTetsuoUmezawa(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tetsuo Umezawa")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARCHER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
