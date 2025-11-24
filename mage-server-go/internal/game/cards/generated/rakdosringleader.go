package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rakdos Ringleader", NewRakdosRingleader)
}

// NewRakdosRingleader creates a Rakdos Ringleader
// {4}{B}{R} - CREATURE
// FirstStrike
func NewRakdosRingleader(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rakdos Ringleader")
	card.ManaCost = "{4}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SKELETON", "WARRIOR"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1, true)
	// card.AddAbility(ability1)
	return card, nil
}
