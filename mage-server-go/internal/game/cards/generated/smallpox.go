package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Smallpox", NewSmallpox)
}

// NewSmallpox creates a Smallpox
// {B}{B} - SORCERY
func NewSmallpox(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Smallpox")
	card.ManaCost = "{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(1, StaticFilters.FILTER_CONTROLLED_CREATURE)
	//   - DiscardEachPlayerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
