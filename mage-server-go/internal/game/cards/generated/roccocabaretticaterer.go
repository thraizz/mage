package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rocco Cabaretti Caterer", NewRoccoCabarettiCaterer)
}

// NewRoccoCabarettiCaterer creates a Rocco Cabaretti Caterer
// {X}{R}{G}{W} - CREATURE
func NewRoccoCabarettiCaterer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rocco Cabaretti Caterer")
	card.ManaCost = "{X}{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "DRUID"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSearchLibraryPutInPlayEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
