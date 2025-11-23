package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Irenicuss Vile Duplication", NewIrenicussVileDuplication)
}

// NewIrenicussVileDuplication creates a Irenicuss Vile Duplication
// {3}{U} - SORCERY
// Flying
func NewIrenicussVileDuplication(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Irenicuss Vile Duplication")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	// card.AddAbility(ability1)
	return card, nil
}
