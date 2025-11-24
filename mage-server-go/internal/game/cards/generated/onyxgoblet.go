package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Onyx Goblet", NewOnyxGoblet)
}

// NewOnyxGoblet creates a Onyx Goblet
// {2}{B} - ARTIFACT
func NewOnyxGoblet(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Onyx Goblet")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewLoseLifeEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}