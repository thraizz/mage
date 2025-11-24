package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blood For The Blood God", NewBloodForTheBloodGod)
}

// NewBloodForTheBloodGod creates a Blood For The Blood God
// {8}{B}{B}{R} - INSTANT
func NewBloodForTheBloodGod(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blood For The Blood God")
	card.ManaCost = "{8}{B}{B}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardHandControllerEffect()
	// card.AddAbility(ability0)
	return card, nil
}