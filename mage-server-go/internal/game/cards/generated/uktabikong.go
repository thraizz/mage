package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Uktabi Kong", NewUktabiKong)
}

// NewUktabiKong creates a Uktabi Kong
// {5}{G}{G}{G} - CREATURE
// Trample
func NewUktabiKong(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Uktabi Kong")
	card.ManaCost = "{5}{G}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"APE"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
