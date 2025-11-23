package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Mordor On The March", NewMordorOnTheMarch)
}

// NewMordorOnTheMarch creates a Mordor On The March
// {3}{B}{R} - SORCERY
func NewMordorOnTheMarch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mordor On The March")
	card.ManaCost = "{3}{B}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(source.getControllerId(), null, false)
	// card.AddAbility(ability0)
	return card, nil
}
