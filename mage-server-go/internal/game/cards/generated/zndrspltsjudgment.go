package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Zndrsplts Judgment", NewZndrspltsJudgment)
}

// NewZndrspltsJudgment creates a Zndrsplts Judgment
// {4}{U} - SORCERY
func NewZndrspltsJudgment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Zndrsplts Judgment")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(player.getId())
	// card.AddAbility(ability0)
	return card, nil
}
