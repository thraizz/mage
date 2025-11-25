package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Convert To Slime", NewConvertToSlime)
}

// NewConvertToSlime creates a Convert To Slime
// {3}{B}{G} - SORCERY
func NewConvertToSlime(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Convert To Slime")
	card.ManaCost = "{3}{B}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddTarget(abilities.NewTargetRequirement(0, 1, abilities.NewArtifactTargetFilter())).
		AddTarget(abilities.NewTargetRequirement(0, 1, abilities.NewCreatureTargetFilter())).
		AddTarget(abilities.NewTargetRequirement(0, 1, abilities.NewEnchantmentTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
