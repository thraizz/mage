package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elmar Ulvenwald Informant", NewElmarUlvenwaldInformant)
}

// NewElmarUlvenwaldInformant creates a Elmar Ulvenwald Informant
// {1}{R}{G} - CREATURE
// Haste
func NewElmarUlvenwaldInformant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elmar Ulvenwald Informant")
	card.ManaCost = "{1}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewUntapEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
