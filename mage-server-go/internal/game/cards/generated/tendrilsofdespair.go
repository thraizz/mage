package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tendrils Of Despair", NewTendrilsOfDespair)
}

// NewTendrilsOfDespair creates a Tendrils Of Despair
// {B} - SORCERY
func NewTendrilsOfDespair(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tendrils Of Despair")
	card.ManaCost = "{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(2)
	//
	// Targets:
	//   - abilities.NewOpponentTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
