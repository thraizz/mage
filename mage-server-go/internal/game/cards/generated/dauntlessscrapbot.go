package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dauntless Scrapbot", NewDauntlessScrapbot)
}

// NewDauntlessScrapbot creates a Dauntless Scrapbot
// {3} - ARTIFACT CREATURE
func NewDauntlessScrapbot(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dauntless Scrapbot")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: ExileGraveyardAllPlayersEffect(StaticFilters.FILTER_CARD, TargetController.OPPONE...)
	// card.AddAbility(ability0)
	return card, nil
}
