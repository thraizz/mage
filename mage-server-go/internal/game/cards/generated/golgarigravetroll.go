package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Golgari Grave Troll", NewGolgariGraveTroll)
}

// NewGolgariGraveTroll creates a Golgari Grave Troll
// {4}{G} - CREATURE
func NewGolgariGraveTroll(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Golgari Grave Troll")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TROLL", "SKELETON"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RegenerateSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
