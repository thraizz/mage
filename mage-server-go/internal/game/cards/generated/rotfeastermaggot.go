package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rotfeaster Maggot", NewRotfeasterMaggot)
}

// NewRotfeasterMaggot creates a Rotfeaster Maggot
// {4}{B} - CREATURE
func NewRotfeasterMaggot(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rotfeaster Maggot")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: RotfeasterMaggotExileEffect()
	// card.AddAbility(ability0)
	return card, nil
}
