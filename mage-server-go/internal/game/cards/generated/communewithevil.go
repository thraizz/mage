package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Commune With Evil", NewCommuneWithEvil)
}

// NewCommuneWithEvil creates a Commune With Evil
// {2}{B} - SORCERY
func NewCommuneWithEvil(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Commune With Evil")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, 1, PutCards.HAND, PutCards.GRA...)
	// card.AddAbility(ability0)
	return card, nil
}
