package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Isleback Spawn", NewIslebackSpawn)
}

// NewIslebackSpawn creates a Isleback Spawn
// {5}{U}{U} - CREATURE
func NewIslebackSpawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Isleback Spawn")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KRAKEN"}
	card.Power = "4"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(4, 8)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
