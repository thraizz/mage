package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Crowded Crypt", NewCrowdedCrypt)
}

// NewCrowdedCrypt creates a Crowded Crypt
// {2}{B} - ARTIFACT
func NewCrowdedCrypt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crowded Crypt")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("ZombieDecayedToken")
	if err != nil {
		return nil, err
	}
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewCreateTokenEffect(token1_0, xValue)).
		Build()
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.NewCounter("corpse", 1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
