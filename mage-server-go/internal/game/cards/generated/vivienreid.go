package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vivien Reid", NewVivienReid)
}

// NewVivienReid creates a Vivien Reid
// {3}{G}{G} - PLANESWALKER
func NewVivienReid(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vivien Reid")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"VIVIEN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, 1, StaticFilters.FILTER_CARD_C...)
	// card.AddAbility(ability1)
	return card, nil
}
