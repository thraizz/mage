package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Zhalfirin Crusader", NewZhalfirinCrusader)
}

// NewZhalfirinCrusader creates a Zhalfirin Crusader
// {1}{W}{W} - CREATURE
func NewZhalfirinCrusader(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Zhalfirin Crusader")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RedirectDamageFromSourceToTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
