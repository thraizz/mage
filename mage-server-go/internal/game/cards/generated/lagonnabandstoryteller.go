package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lagonna Band Storyteller", NewLagonnaBandStoryteller)
}

// NewLagonnaBandStoryteller creates a Lagonna Band Storyteller
// {3}{W} - CREATURE
func NewLagonnaBandStoryteller(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lagonna Band Storyteller")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CENTAUR", "ADVISOR"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: LagonnaBandStorytellerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
