package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Push Pull", NewPushPull)
}

// NewPushPull creates a Push Pull
// {1}{W/B} - SORCERY
func NewPushPull(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Push Pull")
	card.ManaCost = "{1}{W/B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("sacrifice them", controller.getId())
	// card.AddAbility(ability0)
	return card, nil
}
