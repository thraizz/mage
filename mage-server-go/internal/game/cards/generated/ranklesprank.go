package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rankles Prank", NewRanklesPrank)
}

// NewRanklesPrank creates a Rankles Prank
// {2}{B}{B} - SORCERY
func NewRanklesPrank(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rankles Prank")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(2, filter)
	//   - DiscardEachPlayerEffect(2, false)
	// card.AddAbility(ability0)
	return card, nil
}
