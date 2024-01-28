package uri_builder

import (
	"strings"
	"testing"
)

func TestUrlBuilder(t *testing.T) {
	builder := UrlBuilder{}

	// Test setting SSL
	builder.Ssl(true)
	if builder.GetSsl() != true {
		t.Errorf("Expected SSL to be true, got false")
	}
	builder.Ssl(false)
	if builder.GetSsl() != false {
		t.Errorf("Expected SSL to be false, got true")
	}

	// Test setting URI
	uri := "example.com"
	builder.Uri(uri)
	if builder.GetUri() != uri {
		t.Errorf("Expected URI to be %s, got %s", uri, builder.GetUri())
	}

	uri = "example.com:3000/me"
	builder.Uri(uri)
	if builder.GetUri() != uri {
		t.Errorf("Expected URI to be %s, got %s", uri, builder.GetUri())
	}

	// Test building URI with SSL
	uri = "example.com"
	builder.Ssl(true)
	builder.Uri(uri)
	expectedURI := "https://" + uri
	_, builtURI := builder.BuildURI()
	if builtURI != expectedURI {
		t.Errorf("Expected built URI to be %s, got %s", expectedURI, builtURI)
	}

	// Test building URI without SSL
	uri = "example.com:3000/me"
	builder.Ssl(false)
	builder.Uri(uri)
	expectedURI = "http://" + uri
	_, builtURI = builder.BuildURI()
	if builtURI != expectedURI {
		t.Errorf("Expected built URI to be %s, got %s", expectedURI, builtURI)
	}

	// Test building URI with existing scheme
	uri = "https://example.com:3000/me"
	builder.Uri(uri)
	builder.Ssl(false)
	_, builtURI = builder.BuildURI()
	if !strings.HasPrefix(builtURI, "https://") {
		t.Errorf("Expected built URI to start with https://, got %s", builtURI)
	}
	if builtURI != uri {
		t.Errorf("Expected built URI to be %s, got %s", expectedURI, builtURI)
	}
}
