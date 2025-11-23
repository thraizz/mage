package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Emerge From The Cocoon", NewEmergeFromTheCocoon)
}

// NewEmergeFromTheCocoon creates a Emerge From The Cocoon
// {4}{W} - SORCERY
func NewEmergeFromTheCocoon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Emerge From The Cocoon")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
