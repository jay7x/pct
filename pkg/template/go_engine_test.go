package template_test

import (
	"testing"

	"github.com/jay7x/pct/pkg/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoEngineName(t *testing.T) {
	eng, ok := template.Get("go")
	require.True(t, ok)
	assert.Equal(t, "go", eng.Name())
}

func TestGoEngineRender(t *testing.T) {
	eng, ok := template.Get("go")
	require.True(t, ok)

	tests := []struct {
		name    string
		content string
		vars    map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "simple variable substitution",
			content: "Hello {{.name}}!",
			vars:    map[string]interface{}{"name": "world"},
			want:    "Hello world!",
		},
		{
			name:    "toClassName filter",
			content: "class {{.name | toClassName}} { }",
			vars:    map[string]interface{}{"name": "my_class"},
			want:    "class My_class { }",
		},
		{
			name:    "ns2path filter",
			content: "{{.name | ns2path}}",
			vars:    map[string]interface{}{"name": "profile::base"},
			want:    "profile/base",
		},
		{
			name:    "range loop",
			content: "{{range $i, $v := .items}}{{if $i}},{{end}}{{$v}}{{end}}",
			vars:    map[string]interface{}{"items": []string{"a", "b", "c"}},
			want:    "a,b,c",
		},
		{
			name:    "nil vars",
			content: "static text",
			vars:    nil,
			want:    "static text",
		},
		{
			name:    "empty content",
			content: "",
			vars:    map[string]interface{}{},
			want:    "",
		},
		{
			name:    "dot access on nested map",
			content: "{{.puppet_module.version}}",
			vars:    map[string]interface{}{"puppet_module": map[string]interface{}{"version": "0.1.0"}},
			want:    "0.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eng.Render(tt.content, tt.vars)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGoEngineRenderError(t *testing.T) {
	eng, ok := template.Get("go")
	require.True(t, ok)

	_, err := eng.Render("{{.foo", nil)
	assert.Error(t, err)
}

func TestGetUnknown(t *testing.T) {
	_, ok := template.Get("nonexistent")
	assert.False(t, ok)
}
