package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Casal Lurkwood Pathfinder", NewCasalLurkwoodPathfinder)
}

// NewCasalLurkwoodPathfinder creates a Casal Lurkwood Pathfinder
// {3}{G} - CREATURE
// Vigilance
func NewCasalLurkwoodPathfinder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Casal Lurkwood Pathfinder")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TIEFLING", "DRUID"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSearchLibraryPutInPlayEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - TransformSourceEffect()
	//   - DoIfCostPaid(                 new TransformSourceEffect().setTe...)
	// card.AddAbility(ability2)
	return card, nil
}
