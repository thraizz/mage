package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Profane Procession", NewProfaneProcession)
}

// NewProfaneProcession creates a Profane Procession
// {1}{W}{B} - ENCHANTMENT
func NewProfaneProcession(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Profane Procession")
	card.ManaCost = "{1}{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ExileTargetForSourceEffect()
	//   - TransformSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
