package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Innocuous Researcher", NewInnocuousResearcher)
}

// NewInnocuousResearcher creates a Innocuous Researcher
// {3}{G} - CREATURE
func NewInnocuousResearcher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Innocuous Researcher")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CENTAUR", "DETECTIVE"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: UntapAllLandsControllerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
