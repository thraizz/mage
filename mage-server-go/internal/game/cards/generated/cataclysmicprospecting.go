package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cataclysmic Prospecting", NewCataclysmicProspecting)
}

// NewCataclysmicProspecting creates a Cataclysmic Prospecting
// {X}{R}{R} - SORCERY
func NewCataclysmicProspecting(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cataclysmic Prospecting")
	card.ManaCost = "{X}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(GetXValue.instance, StaticFilters.FILTER_PERMANENT...)
	// card.AddAbility(ability0)
	return card, nil
}
