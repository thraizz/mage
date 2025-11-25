package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mizzixs Mastery", NewMizzixsMastery)
}

// NewMizzixsMastery creates a Mizzixs Mastery
// {3}{R} - SORCERY
func NewMizzixsMastery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mizzixs Mastery")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: OverloadAbility
	//   - Effect: MizzixsMasteryOverloadEffect()
	// card.AddAbility(ability0)
	return card, nil
}
