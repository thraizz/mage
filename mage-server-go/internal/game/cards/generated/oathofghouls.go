package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Oath Of Ghouls", NewOathOfGhouls)
}

// NewOathOfGhouls creates a Oath Of Ghouls
// {1}{B} - ENCHANTMENT
func NewOathOfGhouls(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Oath Of Ghouls")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfUpkeepTriggeredAbility
	//   - Effect: OathOfGhoulsEffect()
	// card.AddAbility(ability0)
	return card, nil
}
