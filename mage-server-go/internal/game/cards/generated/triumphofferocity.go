package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Triumph Of Ferocity", NewTriumphOfFerocity)
}

// NewTriumphOfFerocity creates a Triumph Of Ferocity
// {2}{G} - ENCHANTMENT
func NewTriumphOfFerocity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Triumph Of Ferocity")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
