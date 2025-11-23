package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Homarid Spawning Bed", NewHomaridSpawningBed)
}

// NewHomaridSpawningBed creates a Homarid Spawning Bed
// {U}{U} - ENCHANTMENT
func NewHomaridSpawningBed(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Homarid Spawning Bed")
	card.ManaCost = "{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("CamaridToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, SacrificeCostManaValue.CREATURE)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
