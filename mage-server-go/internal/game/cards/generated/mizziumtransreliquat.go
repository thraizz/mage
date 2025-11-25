package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mizzium Transreliquat", NewMizziumTransreliquat)
}

// NewMizziumTransreliquat creates a Mizzium Transreliquat
// {3} - ARTIFACT
func NewMizziumTransreliquat(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mizzium Transreliquat")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - MizziumTransreliquatCopyEffect()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - MizziumTransreliquatCopyAndGainAbilityEffect()
	// card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - MizziumTransreliquatCopyAndGainAbilityEffect()
	// card.AddAbility(ability2)
	return card, nil
}
