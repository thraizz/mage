package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dogged Pursuit", NewDoggedPursuit)
}

// NewDoggedPursuit creates a Dogged Pursuit
// {3}{B} - ENCHANTMENT
func NewDoggedPursuit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dogged Pursuit")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: LoseLifeOpponentsEffect()
	// card.AddAbility(ability0)
	return card, nil
}
