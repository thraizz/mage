package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rals Outburst", NewRalsOutburst)
}

// NewRalsOutburst creates a Rals Outburst
// {2}{U}{R} - INSTANT
func NewRalsOutburst(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rals Outburst")
	card.ManaCost = "{2}{U}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(2, 1, PutCards.HAND, PutCards.GRAVEYARD)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
