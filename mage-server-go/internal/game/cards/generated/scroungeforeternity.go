package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Scrounge For Eternity", NewScroungeForEternity)
}

// NewScroungeForEternity creates a Scrounge For Eternity
// {2}{B} - SORCERY
func NewScroungeForEternity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scrounge For Eternity")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
