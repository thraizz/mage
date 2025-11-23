package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abominable Treefolk", NewAbominableTreefolk)
}

// NewAbominableTreefolk creates a Abominable Treefolk
// {2}{G}{U} - CREATURE
// Trample
func NewAbominableTreefolk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abominable Treefolk")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TREEFOLK"}
	card.Supertypes = []string{"SNOW"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewTapEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
