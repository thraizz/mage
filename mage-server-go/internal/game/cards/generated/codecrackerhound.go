package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Codecracker Hound", NewCodecrackerHound)
}

// NewCodecrackerHound creates a Codecracker Hound
// {2}{U} - CREATURE
func NewCodecrackerHound(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Codecracker Hound")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DOG"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(2, 1, PutCards.HAND, PutCards.GRAVEYARD)
	// card.AddAbility(ability0)
	return card, nil
}
