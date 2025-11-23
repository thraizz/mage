package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Repudiate Replicate", NewRepudiateReplicate)
}

// NewRepudiateReplicate creates a Repudiate Replicate
// {G/U}{G/U} - INSTANT
func NewRepudiateReplicate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Repudiate Replicate")
	card.ManaCost = "{G/U}{G/U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
