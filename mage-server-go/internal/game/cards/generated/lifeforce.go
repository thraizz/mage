package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Lifeforce", NewLifeforce)
}

// NewLifeforce creates a Lifeforce
// {G}{G} - ENCHANTMENT
func NewLifeforce(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lifeforce")
	card.ManaCost = "{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewCounterSpellEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
