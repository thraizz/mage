package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elixir Of Immortality", NewElixirOfImmortality)
}

// NewElixirOfImmortality creates a Elixir Of Immortality
// {1} - ARTIFACT
func NewElixirOfImmortality(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elixir Of Immortality")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ElixerOfImmortalityEffect()
	//
	// Costs:
	//   - AddTapCost()
	//   - AddManaCost("{2}")
	// card.AddAbility(ability0)
	return card, nil
}
