package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Momir Vig Simic Visionary", NewMomirVigSimicVisionary)
}

// NewMomirVigSimicVisionary creates a Momir Vig Simic Visionary
// {3}{G}{U} - CREATURE
func NewMomirVigSimicVisionary(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Momir Vig Simic Visionary")
	card.ManaCost = "{3}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSearchLibraryPutOnTopEffect(abilities.NewTargetRequirement(0, 1, abilities.NewCreatureTargetFilter()), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
