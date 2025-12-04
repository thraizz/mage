package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Conduit Of Worlds", NewConduitOfWorlds)
}

// NewConduitOfWorlds creates a Conduit Of Worlds
// {2}{G}{G} - ARTIFACT
func NewConduitOfWorlds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Conduit Of Worlds")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: ConduitOfWorldsEffect()
	// card.AddAbility(ability0)
	return card, nil
}
