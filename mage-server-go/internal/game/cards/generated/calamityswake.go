package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Calamitys Wake", NewCalamitysWake)
}

// NewCalamitysWake creates a Calamitys Wake
// {1}{W} - INSTANT
func NewCalamitysWake(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Calamitys Wake")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ExileGraveyardAllPlayersEffect()
	// card.AddAbility(ability0)
	return card, nil
}
