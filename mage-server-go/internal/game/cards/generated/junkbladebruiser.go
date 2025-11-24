package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Junkblade Bruiser", NewJunkbladeBruiser)
}

// NewJunkbladeBruiser creates a Junkblade Bruiser
// {3}{R/G}{R/G} - CREATURE
// Trample
func NewJunkbladeBruiser(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Junkblade Bruiser")
	card.ManaCost = "{3}{R/G}{R/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RACCOON", "BERSERKER"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(2, 1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}