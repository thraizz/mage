package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tezzeret Master Of The Bridge", NewTezzeretMasterOfTheBridge)
}

// NewTezzeretMasterOfTheBridge creates a Tezzeret Master Of The Bridge
// {4}{U}{B} - PLANESWALKER
func NewTezzeretMasterOfTheBridge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tezzeret Master Of The Bridge")
	card.ManaCost = "{4}{U}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TEZZERET"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewReturnFromGraveyardToHandTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
