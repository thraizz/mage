package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Together Forever", NewTogetherForever)
}

// NewTogetherForever creates a Together Forever
// {W}{W} - ENCHANTMENT
func NewTogetherForever(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Together Forever")
	card.ManaCost = "{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddEffect(abilities.NewReturnFromGraveyardToHandTargetEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
