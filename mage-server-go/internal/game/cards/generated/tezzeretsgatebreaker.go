package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tezzerets Gatebreaker", NewTezzeretsGatebreaker)
}

// NewTezzeretsGatebreaker creates a Tezzerets Gatebreaker
// {4} - ARTIFACT
func NewTezzeretsGatebreaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tezzerets Gatebreaker")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(5, 1, filter, PutCards.HAND, PutCards.BOTTOM_RANDO...)
	// card.AddAbility(ability0)
	return card, nil
}
