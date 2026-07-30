package pct_config_processor_test

import (
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/jay7x/pct/internal/pkg/pct"
	"github.com/jay7x/pct/pkg/config_processor"
	"github.com/jay7x/pct/pkg/install"

	"github.com/jay7x/pct/internal/pkg/pct_config_processor"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type CheckConfigTest struct {
	name           string
	mockConfigFile bool
	configFilePath string
	configFileYaml string
	errorMsg       string
}

func TestPctConfigProcessor_CheckConfig(t *testing.T) {
	tests := []CheckConfigTest{
		{
			name:           "When config not found",
			mockConfigFile: false,
			configFilePath: "my/missing/pct-config.yml",
			errorMsg:       "file does not exist",
		},
		{
			name:           "When config valid",
			mockConfigFile: true,
			configFilePath: "my/valid/pct-config.yml",

			configFileYaml: `---
template:
  id: test-template
  author: test-user
  version: 0.1.0
`,
			errorMsg: "",
		},
		{
			name:           "When config invalid",
			mockConfigFile: true,
			configFilePath: "my/invalid/pct-config.yml",
			// This is invalid because it starts with tabs which the parses errors on
			configFileYaml: `---
			foo: bar
			`,
			errorMsg: "parsing config: yaml",
		},
		{
			name:           "When config missing author",
			mockConfigFile: true,
			configFilePath: "my/missing/author/pct-config.yml",

			configFileYaml: `---
template:
  id: test-template
  version: 0.1.0
`,
			errorMsg: `The following attributes are missing in .+:\s+\* author`,
		},
		{
			name:           "When config missing id",
			mockConfigFile: true,
			configFilePath: "my/missing/id/pct-config.yml",

			configFileYaml: `---
template:
  author: test-user
  version: 0.1.0
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id`,
		},
		{
			name:           "When config missing version",
			mockConfigFile: true,
			configFilePath: "my/missing/version/pct-config.yml",

			configFileYaml: `---
template:
  author: test-user
  id: test-template
`,
			errorMsg: `The following attributes are missing in .+:\s+\* version`,
		},
		{
			name:           "When config missing author, id, and version",
			mockConfigFile: true,
			configFilePath: "my/missing/version/pct-config.yml",

			configFileYaml: `---
template:
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id\s+\* author\s+\* version`,
		},
		{
			name:           "When config missing template key",
			mockConfigFile: true,
			configFilePath: "my/missing/version/pct-config.yml",

			configFileYaml: `---
foo: bar
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id\s+\* author\s+\* version`,
		},
		{
			name:           "When config author starts with '.'",
			mockConfigFile: true,
			configFilePath: "my/dotfile/author/pct-config.yml",

			configFileYaml: `---
template:
  id: test-template
  author: .hidden
  version: 0.1.0
`,
			errorMsg: `The following attributes are missing in .+:\s+\* author must not start with '\.'`,
		},
		{
			name:           "When config id starts with '.'",
			mockConfigFile: true,
			configFilePath: "my/dotfile/id/pct-config.yml",

			configFileYaml: `---
template:
  id: .hidden-template
  author: test-user
  version: 0.1.0
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id must not start with '\.'`,
		},
		{
			name:           "When config author and id both start with '.'",
			mockConfigFile: true,
			configFilePath: "my/dotfile/both/pct-config.yml",

			configFileYaml: `---
template:
  id: .hidden-template
  author: .hidden
  version: 0.1.0
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id must not start with '\.'\s+\* author must not start with '\.'`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}

			if tt.mockConfigFile {
				dir := filepath.Dir(tt.configFilePath)
				_ = afs.Mkdir(dir, 0o0750)
				config, _ := afs.Create(tt.configFilePath)
				_, _ = config.WriteString(tt.configFileYaml)
			}

			configProcessor := pct_config_processor.PctConfigProcessor{AFS: afs}

			err := configProcessor.CheckConfig(tt.configFilePath)

			if tt.errorMsg != "" && err != nil {
				assert.Regexp(t, regexp.MustCompile(tt.errorMsg), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPctConfigProcessor_GetConfigMetadata(t *testing.T) {
	type args struct {
		configFile string
	}
	configParentPath := "path/to/extract/to/"

	tests := []struct {
		name           string
		args           args
		wantMetadata   config_processor.ConfigMetadata
		wantErr        bool
		templateConfig string // Leave blank for config file not to be created
	}{
		{
			name: "Successfully gets config metadata",
			args: args{
				configFile: filepath.Join(configParentPath, "pct-config.yml"),
			},
			wantMetadata: config_processor.ConfigMetadata{Author: "test-user", Id: "full-project", Version: "0.1.0"},
			templateConfig: `---
template:
  id: full-project
  author: test-user
  version: 0.1.0
`,
		},
		{
			name: "Missing vital metadata from pct-config.yml (id omitted)",
			args: args{
				configFile: filepath.Join(configParentPath, "pct-config.yml"),
			},
			wantErr:      true,
			wantMetadata: config_processor.ConfigMetadata{},
			templateConfig: `---
template:
  author: test-user
  version: 0.1.0
`,
		},
		{
			name: "Malformed pct-config (extra indentation)",
			args: args{
				configFile: filepath.Join(configParentPath, "pct-config.yml"),
			},
			wantErr:      true,
			wantMetadata: config_processor.ConfigMetadata{},
			templateConfig: `---
	template:
		id: full-project
  	author: test-user
  	version: 0.1.0
`, // Contains an erroneous extra indent
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Instantiate afs
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			p := &pct_config_processor.PctConfigProcessor{
				AFS: afs,
			}

			// Create all useful directories
			_ = afs.MkdirAll(configParentPath, 0o0750)
			if tt.templateConfig != "" {
				config, _ := afs.Create(tt.args.configFile)
				_, _ = config.WriteString(tt.templateConfig)
			}

			gotMetadata, err := p.GetConfigMetadata(tt.args.configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMetadata, tt.wantMetadata) {
				t.Errorf("GetConfigMetadata() gotMetadata = %v, want %v", gotMetadata, tt.wantMetadata)
			}
		})
	}
}

func TestPctConfigProcessor_ReadConfig(t *testing.T) {
	type fields struct {
		AFS *afero.Afero
	}
	type args struct {
		configFile string
	}
	configParentPath := "path/to/extract/to/"

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantInfo       pct.PuppetContentTemplateInfo
		wantErr        bool
		templateConfig string
	}{
		{
			name: "Successfully read file and return filled struct",
			wantInfo: pct.PuppetContentTemplateInfo{
				Template: pct.PuppetContentTemplate{
					ConfigParams: install.ConfigParams{
						Id:      "full-project",
						Author:  "test-user",
						Version: "0.1.0",
					},
					Type:    "class",
					Display: "Full Project",
					URL:     "github.com/puppetlabs/pct",
				},
			},
			templateConfig: `---
template:
  id: full-project
  author: test-user
  version: 0.1.0
  type: class
  display: Full Project
  url: github.com/puppetlabs/pct
`,
		},
		{
			name:    "Returns an error because of the malformed template config",
			wantErr: true,
			templateConfig: `---
	template:
		id: full-project
		author: test-user
		version: 0.1.0
		type: class
		display: Full Project
		url: github.com/puppetlabs/pct
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Instantiate afs
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			p := &pct_config_processor.PctConfigProcessor{
				AFS: afs,
			}

			// Create all useful directories
			_ = afs.MkdirAll(configParentPath, 0o0750)
			if tt.templateConfig != "" {
				config, _ := afs.Create(tt.args.configFile)
				_, _ = config.WriteString(tt.templateConfig)
			}

			gotInfo, err := p.ReadConfig(tt.args.configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("ReadConfig() gotInfo = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}
