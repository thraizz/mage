package server

import "testing"

func TestParseCardLine_StripsMTGOStyleSuffix(t *testing.T) {
	qty, name := parseCardLine("1x Beast Within (eoc) 93")
	if qty != 1 {
		t.Fatalf("expected qty=1, got %d", qty)
	}
	if name != "Beast Within" {
		t.Fatalf("expected name=%q, got %q", "Beast Within", name)
	}
}

func TestParseDeckListText_StripsSuffixesInAllSections(t *testing.T) {
	input := "Commander:\n1x Atraxa, Praetors' Voice (cmm) 350\n\n4x Beast Within (eoc) 93\n\nSideboard:\n2 Nature's Claim (m11) 193\n"
	main, side, commanders := parseDeckListText(input)

	if len(commanders) != 1 || commanders[0] != "Atraxa, Praetors' Voice" {
		t.Fatalf("unexpected commanders: %#v", commanders)
	}
	if len(main) != 4 {
		t.Fatalf("expected 4 main deck cards, got %d (%#v)", len(main), main)
	}
	for _, c := range main {
		if c != "Beast Within" {
			t.Fatalf("expected main card %q, got %q", "Beast Within", c)
		}
	}
	if len(side) != 2 {
		t.Fatalf("expected 2 sideboard cards, got %d (%#v)", len(side), side)
	}
	for _, c := range side {
		if c != "Nature's Claim" {
			t.Fatalf("expected sideboard card %q, got %q", "Nature's Claim", c)
		}
	}
}

func TestNormalizeImportedCardName_DoesNotStripAdventureDoubleSlash(t *testing.T) {
	// We only strip "(SET) 123" suffixes here; adventure handling is DB-backed in resolveCardNames.
	in := "Brazen Borrower // Petty Theft (eld) 39"
	out := normalizeImportedCardName(in)
	if out != "Brazen Borrower // Petty Theft" {
		t.Fatalf("expected %q, got %q", "Brazen Borrower // Petty Theft", out)
	}
}

func TestParseCardLine_PreservesAdventureSeparatorSpacing(t *testing.T) {
	// Some exporters omit spaces around the separator; the server should still handle it later.
	qty, name := parseCardLine("1x Brazen Borrower//Petty Theft (eld) 39")
	if qty != 1 {
		t.Fatalf("expected qty=1, got %d", qty)
	}
	if name != "Brazen Borrower//Petty Theft" {
		t.Fatalf("expected name=%q, got %q", "Brazen Borrower//Petty Theft", name)
	}
}
