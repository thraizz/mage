package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Broodguard Elite", NewBroodguardElite)
}

// NewBroodguardElite creates a Broodguard Elite
// {X}{G}{G} - CREATURE
func NewBroodguardElite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Broodguard Elite")
	card.ManaCost = "{X}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT", "KNIGHT"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LeavesBattlefieldTriggeredAbility
	//   - Effect: PutSourceCountersOnTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
