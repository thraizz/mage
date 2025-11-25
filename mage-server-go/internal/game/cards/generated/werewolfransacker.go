package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Werewolf Ransacker", NewWerewolfRansacker)
}

// NewWerewolfRansacker creates a Werewolf Ransacker
//   - CREATURE
func NewWerewolfRansacker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Werewolf Ransacker")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WEREWOLF"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: TransformIntoSourceTriggeredAbility
	//   - Effect: WerewolfRansackerEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewArtifactTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
