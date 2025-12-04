package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gaeas Touch", NewGaeasTouch)
}

// NewGaeasTouch creates a Gaeas Touch
// {G}{G} - ENCHANTMENT
func NewGaeasTouch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gaeas Touch")
	card.ManaCost = "{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LimitedTimesPerTurnActivatedAbility
	//   - Effect: PutCardFromHandOntoBattlefieldEffect()
	// card.AddAbility(ability0)
	return card, nil
}
