package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wydwen The Biting Gale", NewWydwenTheBitingGale)
}

// NewWydwenTheBitingGale creates a Wydwen The Biting Gale
// {2}{U}{B} - CREATURE
// Flash, Flying
func NewWydwenTheBitingGale(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wydwen The Biting Gale")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FAERIE", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewReturnToHandSourceEffect()).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
