package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Urzas Ruinous Blast", NewUrzasRuinousBlast)
}

// NewUrzasRuinousBlast creates a Urzas Ruinous Blast
// {4}{W} - SORCERY
func NewUrzasRuinousBlast(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Urzas Ruinous Blast")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"SORCERY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewExileAllEffect(abilities.NewAnyTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
