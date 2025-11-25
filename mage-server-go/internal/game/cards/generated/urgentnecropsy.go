package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Urgent Necropsy", NewUrgentNecropsy)
}

// NewUrgentNecropsy creates a Urgent Necropsy
// {2}{B}{G} - INSTANT
func NewUrgentNecropsy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Urgent Necropsy")
	card.ManaCost = "{2}{B}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddTargets(0, 1, abilities.NewArtifactTargetFilter()).
		AddTargets(0, 1, abilities.NewCreatureTargetFilter()).
		AddTargets(0, 1, abilities.NewEnchantmentTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
