package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ayara First Of Locthwain", NewAyaraFirstOfLocthwain)
}

// NewAyaraFirstOfLocthwain creates a Ayara First Of Locthwain
// {B}{B}{B} - CREATURE
func NewAyaraFirstOfLocthwain(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ayara First Of Locthwain")
	card.ManaCost = "{B}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
