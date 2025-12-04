package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Red Suns Twilight", NewRedSunsTwilight)
}

// NewRedSunsTwilight creates a Red Suns Twilight
// {X}{R}{R} - SORCERY
func NewRedSunsTwilight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Red Suns Twilight")
	card.ManaCost = "{X}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(player.getId(), null, true)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(0, 1, abilities.NewArtifactTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
