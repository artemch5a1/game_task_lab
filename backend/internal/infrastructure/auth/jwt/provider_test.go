package jwt

import (
	"testing"
	"time"

	specifictype "example/web-service-gin/internal/domain/specific_type"

	"github.com/google/uuid"
)

func TestProvider_IssueAndParse(t *testing.T) {
	p := NewProvider("test-secret", "test-issuer", 1*time.Hour)
	uid := uuid.New()

	token, err := p.Issue(nil, uid, specifictype.RoleAdmin)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	claims, err := p.Parse(token)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if claims.Subject != uid.String() {
		t.Fatalf("expected sub %s, got %s", uid.String(), claims.Subject)
	}
	if claims.Role != string(specifictype.RoleAdmin) {
		t.Fatalf("expected role %s, got %s", specifictype.RoleAdmin, claims.Role)
	}
	if claims.Issuer != "test-issuer" {
		t.Fatalf("expected issuer test-issuer, got %s", claims.Issuer)
	}
}

