package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Machine Gods Effigy", NewMachineGodsEffigy)
}

// NewMachineGodsEffigy creates a Machine Gods Effigy
// {4} - ARTIFACT
func NewMachineGodsEffigy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Machine Gods Effigy")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(StaticFilters.FILTER_PERMANENT_CREATURE, new Machi...)
	// card.AddAbility(ability2)
	return card, nil
}
