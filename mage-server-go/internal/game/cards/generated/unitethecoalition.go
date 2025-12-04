package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unite The Coalition", NewUniteTheCoalition)
}

// NewUniteTheCoalition creates a Unite The Coalition
// {2}{W}{U}{B}{R}{G} - INSTANT
func NewUniteTheCoalition(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unite The Coalition")
	card.ManaCost = "{2}{W}{U}{B}{R}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - PhaseOutTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewArtifactOrEnchantmentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
