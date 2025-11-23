package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Eight And A Half Tails", NewEightAndAHalfTails)
}

// NewEightAndAHalfTails creates a Eight And A Half Tails
// {W}{W} - CREATURE
func NewEightAndAHalfTails(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eight And A Half Tails")
	card.ManaCost = "{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FOX", "CLERIC"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewGrantAbilityEffect(ProtectionAbility.from(ObjectColor.WHITE))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
