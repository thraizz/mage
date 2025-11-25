package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chaotic Transformation", NewChaoticTransformation)
}

// NewChaoticTransformation creates a Chaotic Transformation
// {5}{R} - SORCERY
func NewChaoticTransformation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chaotic Transformation")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddTarget(abilities.NewTargetRequirement(0, 1, abilities.NewArtifactTargetFilter())).
		AddTarget(abilities.NewTargetRequirement(0, 1, abilities.NewCreatureTargetFilter())).
		AddTarget(abilities.NewTargetRequirement(0, 1, abilities.NewEnchantmentTargetFilter())).
		AddTarget(abilities.NewTargetRequirement(0, 1, abilities.NewLandTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
