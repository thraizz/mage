package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Devastating Onslaught", NewDevastatingOnslaught)
}

// NewDevastatingOnslaught creates a Devastating Onslaught
// {X}{X}{R} - SORCERY
func NewDevastatingOnslaught(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Devastating Onslaught")
	card.ManaCost = "{X}{X}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(null, null, true, xValue)
	// card.AddAbility(ability0)
	return card, nil
}
