package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Live Or Die", NewLiveOrDie)
}

// NewLiveOrDie creates a Live Or Die
// {3}{B}{B} - INSTANT
func NewLiveOrDie(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Live Or Die")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	//
	// Targets:
	//   - abilities.NewCreatureTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
