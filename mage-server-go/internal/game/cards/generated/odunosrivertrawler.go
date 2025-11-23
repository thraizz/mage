package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Odunos River Trawler", NewOdunosRiverTrawler)
}

// NewOdunosRiverTrawler creates a Odunos River Trawler
// {2}{B} - CREATURE
func NewOdunosRiverTrawler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Odunos River Trawler")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewReturnFromGraveyardToHandTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewReturnFromGraveyardToHandTargetEffect()).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
