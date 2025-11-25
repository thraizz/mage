package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Steelbane Hydra", NewSteelbaneHydra)
}

// NewSteelbaneHydra creates a Steelbane Hydra
// {X}{G}{G} - CREATURE
func NewSteelbaneHydra(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Steelbane Hydra")
	card.ManaCost = "{X}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TURTLE", "HYDRA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
