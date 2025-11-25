package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Cogwork Assembler", NewCogworkAssembler)
}

// NewCogworkAssembler creates a Cogwork Assembler
// {3} - ARTIFACT CREATURE
func NewCogworkAssembler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cogwork Assembler")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ASSEMBLY_WORKER"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(null, CardType.ARTIFACT, true)
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - CogworkAssemblerCreateTokenEffect()
	//
	// Costs:
	//   - AddManaCost("{7}")
	// card.AddAbility(ability1)
	return card, nil
}
