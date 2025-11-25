package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Saheelis Artistry", NewSaheelisArtistry)
}

// NewSaheelisArtistry creates a Saheelis Artistry
// {4}{U}{U} - SORCERY
func NewSaheelisArtistry(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Saheelis Artistry")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	//   - CreateTokenCopyTargetEffect()
	//   - CreateTokenCopyTargetEffect()
	//   - CreateTokenCopyTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewArtifactTargetFilter())
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
