package permissions

import "testing"

func TestDefaultForRoleProfessionalUsesOwnScope(t *testing.T) {
	t.Parallel()

	value := DefaultForRole("professional")
	if value.UserDirectory != ScopeOwn ||
		value.Patients != ScopeOwn ||
		value.Services != ScopeOwn ||
		value.Calendar != ScopeOwn ||
		value.Finance != ScopeOwn ||
		value.AIHistory != ScopeOwn ||
		value.Plans != ScopeOwn {
		t.Fatalf("expected professional to default to own scope, got %#v", value)
	}
}

func TestDefaultForRoleAdminUsesAllScope(t *testing.T) {
	t.Parallel()

	value := DefaultForRole("admin")
	if value.UserDirectory != ScopeAll ||
		value.Patients != ScopeAll ||
		value.Services != ScopeAll ||
		value.Calendar != ScopeAll ||
		value.Finance != ScopeAll ||
		value.AIHistory != ScopeAll ||
		value.Plans != ScopeAll {
		t.Fatalf("expected admin to default to all scope, got %#v", value)
	}
}
