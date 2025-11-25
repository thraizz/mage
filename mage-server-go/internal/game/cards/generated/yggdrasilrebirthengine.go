package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yggdrasil Rebirth Engine", NewYggdrasilRebirthEngine)
}

// NewYggdrasilRebirthEngine creates a Yggdrasil Rebirth Engine
// {3} - ARTIFACT
func NewYggdrasilRebirthEngine(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yggdrasil Rebirth Engine")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: YggdrasilRebirthEngineReturnCreatureEffect()
	// card.AddAbility(ability0)
	return card, nil
}
