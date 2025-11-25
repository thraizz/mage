package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Muse Vessel", NewMuseVessel)
}

// NewMuseVessel creates a Muse Vessel
// {4} - ARTIFACT
func NewMuseVessel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Muse Vessel")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: MuseVesselExileEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - OneShotNonTargetEffect()
	// card.AddAbility(ability1)
	return card, nil
}
