package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sire Of Stagnation", NewSireOfStagnation)
}

// NewSireOfStagnation creates a Sire Of Stagnation
// {4}{U}{B} - CREATURE
func NewSireOfStagnation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sire Of Stagnation")
	card.ManaCost = "{4}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Power = "5"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldOpponentTriggeredAbility
	//   - Effect: ExileCardsFromTopOfLibraryTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
