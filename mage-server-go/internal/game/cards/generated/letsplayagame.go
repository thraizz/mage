package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lets Play A Game", NewLetsPlayAGame)
}

// NewLetsPlayAGame creates a Lets Play A Game
// {3}{B} - SORCERY
func NewLetsPlayAGame(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lets Play A Game")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(                 StaticValue.get(2), false, Target...)
	// card.AddAbility(ability0)
	return card, nil
}
