package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Supreme Leader Snoke", NewSupremeLeaderSnoke)
}

// NewSupremeLeaderSnoke creates a Supreme Leader Snoke
// {U}{B}{R} - PLANESWALKER
func NewSupremeLeaderSnoke(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Supreme Leader Snoke")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SNOKE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewGainControlTargetEffect(abilities.DurationWhileOnBattlefield)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}