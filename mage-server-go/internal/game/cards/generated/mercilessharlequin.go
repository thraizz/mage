package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Merciless Harlequin", NewMercilessHarlequin)
}

// NewMercilessHarlequin creates a Merciless Harlequin
// {2}{B} - CREATURE
func NewMercilessHarlequin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Merciless Harlequin")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ASSASSIN"}
	card.Power = "2"
	card.Toughness = "1"
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
