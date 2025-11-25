package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gontis Aether Heart", NewGontisAetherHeart)
}

// NewGontisAetherHeart creates a Gontis Aether Heart
// {6} - ARTIFACT
func NewGontisAetherHeart(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gontis Aether Heart")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - AddExtraTurnControllerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
