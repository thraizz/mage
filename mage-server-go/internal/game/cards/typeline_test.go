package cards

import (
	"reflect"
	"testing"
)

func TestParseTypeLine(t *testing.T) {
	tests := []struct {
		name     string
		typeLine string
		want     ParsedTypeLine
	}{
		{
			name:     "Simple creature",
			typeLine: "Creature — Human Wizard",
			want: ParsedTypeLine{
				Supertypes: []string{},
				Types:      []string{"CREATURE"},
				Subtypes:   []string{"HUMAN", "WIZARD"},
			},
		},
		{
			name:     "Legendary creature",
			typeLine: "Legendary Creature — Human Wizard",
			want: ParsedTypeLine{
				Supertypes: []string{"LEGENDARY"},
				Types:      []string{"CREATURE"},
				Subtypes:   []string{"HUMAN", "WIZARD"},
			},
		},
		{
			name:     "Basic land",
			typeLine: "Basic Land — Forest",
			want: ParsedTypeLine{
				Supertypes: []string{"BASIC"},
				Types:      []string{"LAND"},
				Subtypes:   []string{"FOREST"},
			},
		},
		{
			name:     "Instant spell",
			typeLine: "Instant",
			want: ParsedTypeLine{
				Supertypes: []string{},
				Types:      []string{"INSTANT"},
				Subtypes:   []string{},
			},
		},
		{
			name:     "Artifact creature",
			typeLine: "Artifact Creature — Golem",
			want: ParsedTypeLine{
				Supertypes: []string{},
				Types:      []string{"ARTIFACT", "CREATURE"},
				Subtypes:   []string{"GOLEM"},
			},
		},
		{
			name:     "Legendary artifact",
			typeLine: "Legendary Artifact",
			want: ParsedTypeLine{
				Supertypes: []string{"LEGENDARY"},
				Types:      []string{"ARTIFACT"},
				Subtypes:   []string{},
			},
		},
		{
			name:     "Land without subtypes",
			typeLine: "Land",
			want: ParsedTypeLine{
				Supertypes: []string{},
				Types:      []string{"LAND"},
				Subtypes:   []string{},
			},
		},
		{
			name:     "Snow land",
			typeLine: "Snow Land — Forest",
			want: ParsedTypeLine{
				Supertypes: []string{"SNOW"},
				Types:      []string{"LAND"},
				Subtypes:   []string{"FOREST"},
			},
		},
		{
			name:     "Legendary snow creature",
			typeLine: "Legendary Snow Creature — Treefolk",
			want: ParsedTypeLine{
				Supertypes: []string{"LEGENDARY", "SNOW"},
				Types:      []string{"CREATURE"},
				Subtypes:   []string{"TREEFOLK"},
			},
		},
		{
			name:     "Planeswalker",
			typeLine: "Legendary Planeswalker — Jace",
			want: ParsedTypeLine{
				Supertypes: []string{"LEGENDARY"},
				Types:      []string{"PLANESWALKER"},
				Subtypes:   []string{"JACE"},
			},
		},
		{
			name:     "Battle",
			typeLine: "Battle — Siege",
			want: ParsedTypeLine{
				Supertypes: []string{},
				Types:      []string{"BATTLE"},
				Subtypes:   []string{"SIEGE"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTypeLine(tt.typeLine)

			if !reflect.DeepEqual(got.Supertypes, tt.want.Supertypes) {
				t.Errorf("ParseTypeLine() supertypes = %v, want %v", got.Supertypes, tt.want.Supertypes)
			}
			if !reflect.DeepEqual(got.Types, tt.want.Types) {
				t.Errorf("ParseTypeLine() types = %v, want %v", got.Types, tt.want.Types)
			}
			if !reflect.DeepEqual(got.Subtypes, tt.want.Subtypes) {
				t.Errorf("ParseTypeLine() subtypes = %v, want %v", got.Subtypes, tt.want.Subtypes)
			}
		})
	}
}

func TestParsedTypeLineHelpers(t *testing.T) {
	creature := ParseTypeLine("Legendary Creature — Human Wizard")
	if !creature.IsCreature() {
		t.Error("Expected IsCreature() to be true")
	}
	if !creature.IsLegendary() {
		t.Error("Expected IsLegendary() to be true")
	}
	if !creature.HasSubtype("Human") {
		t.Error("Expected HasSubtype(Human) to be true")
	}
	if !creature.HasSubtype("WIZARD") {
		t.Error("Expected HasSubtype(WIZARD) to be true")
	}

	instant := ParseTypeLine("Instant")
	if !instant.IsInstant() {
		t.Error("Expected IsInstant() to be true")
	}
	if instant.IsCreature() {
		t.Error("Expected IsCreature() to be false")
	}

	land := ParseTypeLine("Basic Land — Forest")
	if !land.IsLand() {
		t.Error("Expected IsLand() to be true")
	}
	if !land.IsBasic() {
		t.Error("Expected IsBasic() to be true")
	}
	if !land.HasSubtype("Forest") {
		t.Error("Expected HasSubtype(Forest) to be true")
	}
}
