package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Teferi Temporal Archmage", NewTeferiTemporalArchmage)
}

// NewTeferiTemporalArchmage creates a Teferi Temporal Archmage
// {4}{U}{U} - PLANESWALKER
func NewTeferiTemporalArchmage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Teferi Temporal Archmage")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TEFERI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewUntapEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(2, 1, PutCards.HAND, PutCards.BOTTOM_ANY)
	// card.AddAbility(ability1)
	return card, nil
}
