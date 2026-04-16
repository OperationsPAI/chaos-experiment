package handler

import "testing"

func TestSystemTypeConstantsIncludeRegisteredSystems(t *testing.T) {
	for _, system := range []SystemType{SystemSockShop, SystemTeaStore} {
		if !system.IsValid() {
			t.Fatalf("expected %s to be a valid handler system type", system)
		}
	}

	systems := GetAllSystemTypes()
	seen := make(map[SystemType]bool, len(systems))
	for _, system := range systems {
		seen[system] = true
	}

	for _, system := range []SystemType{SystemSockShop, SystemTeaStore} {
		if !seen[system] {
			t.Fatalf("expected %s to be returned by GetAllSystemTypes()", system)
		}
	}
}
