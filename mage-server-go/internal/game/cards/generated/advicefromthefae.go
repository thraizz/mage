package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Advice From The Fae", NewAdviceFromTheFae)
}

// NewAdviceFromTheFae creates a Advice From The Fae
// {2/U}{2/U}{2/U} - SORCERY
func NewAdviceFromTheFae(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Advice From The Fae")
	card.ManaCost = "{2/U}{2/U}{2/U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(5, 2, PutCards.HAND, PutCards.BOTTOM_ANY)
	//   - LookLibraryAndPickControllerEffect(5, 1, PutCards.HAND, PutCards.BOTTOM_ANY)
	// card.AddAbility(ability0)
	return card, nil
}