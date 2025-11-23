package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mirris Guile", NewMirrisGuile)
}

// NewMirrisGuile creates a Mirris Guile
// {G} - ENCHANTMENT
func NewMirrisGuile(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mirris Guile")
	card.ManaCost = "{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryControllerEffect(3)
	// card.AddAbility(ability0)
	return card, nil
}
