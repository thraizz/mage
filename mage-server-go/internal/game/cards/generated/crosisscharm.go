package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crosiss Charm", NewCrosissCharm)
}

// NewCrosissCharm creates a Crosiss Charm
// {U}{B}{R} - INSTANT
func NewCrosissCharm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crosiss Charm")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect(true)).
		// TODO: DestroyTargetEffect with complex parameters
		// TODO: ReturnToHandTargetEffect with complex parameters
		AddEffect(abilities.NewDestroyEffect(true)).
		// TODO: DestroyTargetEffect with complex parameters
		AddTarget(abilities.NewPermanentTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
