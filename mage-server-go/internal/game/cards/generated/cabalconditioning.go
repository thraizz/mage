package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cabal Conditioning", NewCabalConditioning)
}

// NewCabalConditioning creates a Cabal Conditioning
// {6}{B} - SORCERY
func NewCabalConditioning(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cabal Conditioning")
	card.ManaCost = "{6}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(GreatestAmongPermanentsValue.MANAVALUE_CONTROLLED_...)
	// card.AddAbility(ability0)
	return card, nil
}
