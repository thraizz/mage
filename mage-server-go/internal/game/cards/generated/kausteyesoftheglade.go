package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaust Eyes Of The Glade", NewKaustEyesOfTheGlade)
}

// NewKaustEyesOfTheGlade creates a Kaust Eyes Of The Glade
// {R/W}{G} - CREATURE
func NewKaustEyesOfTheGlade(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaust Eyes Of The Glade")
	card.ManaCost = "{R/W}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRYAD", "DETECTIVE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
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
