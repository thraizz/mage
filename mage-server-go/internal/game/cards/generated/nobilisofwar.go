package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nobilis Of War", NewNobilisOfWar)
}

// NewNobilisOfWar creates a Nobilis Of War
// {R/W}{R/W}{R/W}{R/W}{R/W} - CREATURE
// Flying
func NewNobilisOfWar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nobilis Of War")
	card.ManaCost = "{R/W}{R/W}{R/W}{R/W}{R/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "AVATAR"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(2, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
