package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Belltower Sphinx", NewBelltowerSphinx)
}

// NewBelltowerSphinx creates a Belltower Sphinx
// {4}{U} - CREATURE
// Flying
func NewBelltowerSphinx(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Belltower Sphinx")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPHINX"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewMillCardsTargetEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
