package util

import "testing"

func TestToHypenUUID(t *testing.T) {
	uuid := "069a79f444e94726a5befca90e38aaf5"
	expected := "069a79f4-44e9-4726-a5be-fca90e38aaf5"
	if hyphens := ToHypenUUID(uuid); expected != hyphens {
		t.Error("Expected " + expected + " and got: " + hyphens)
	}
}

func TestIsValidUsername(t *testing.T) {
	valid := "Nathanael"
	if !IsValidUsername(valid) {
		t.Error(valid + " is a valid username")
	}
	tooLong := "MichiganDu75StreetFighterTeamKenzafarah"
	if IsValidUsername(tooLong) {
		t.Error(tooLong + " is not a valid username")
	}
	wrongChars := "Yes..."
	if IsValidUsername(wrongChars) {
		t.Error(wrongChars + " is not a valid username")
	}
}
