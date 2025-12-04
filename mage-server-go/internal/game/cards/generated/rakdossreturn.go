package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rakdoss Return", NewRakdossReturn)
}

// NewRakdossReturn creates a Rakdoss Return
// {X}{B}{R} - SORCERY
func NewRakdossReturn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rakdoss Return")
	card.ManaCost = "{X}{B}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(GetXValue.instance)
	// card.AddAbility(ability0)
	return card, nil
}
