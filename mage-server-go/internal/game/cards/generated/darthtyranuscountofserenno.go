package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Darth Tyranus Count Of Serenno", NewDarthTyranusCountOfSerenno)
}

// NewDarthTyranusCountOfSerenno creates a Darth Tyranus Count Of Serenno
// {1}{W}{U}{B} - PLANESWALKER
func NewDarthTyranusCountOfSerenno(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Darth Tyranus Count Of Serenno")
	card.ManaCost = "{1}{W}{U}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DOOKU"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(-6, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}