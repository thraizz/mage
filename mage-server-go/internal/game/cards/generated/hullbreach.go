package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hull Breach", NewHullBreach)
}

// NewHullBreach creates a Hull Breach
// {R}{G} - SORCERY
func NewHullBreach(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hull Breach")
	card.ManaCost = "{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		AddTarget(abilities.NewArtifactTargetFilter()).
		AddTarget(abilities.NewEnchantmentTargetFilter()).
		AddTarget(abilities.NewArtifactTargetFilter()).
		AddTarget(abilities.NewEnchantmentTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
