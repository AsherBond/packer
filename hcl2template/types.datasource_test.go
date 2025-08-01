// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package hcl2template

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/packer/builder/null"
	"github.com/hashicorp/packer/packer"
)

func TestParse_datasource(t *testing.T) {
	defaultParser := getBasicParser()

	tests := []parseTest{
		{"two basic datasources",
			defaultParser,
			parseTestArgs{"testdata/datasources/basic.pkr.hcl", nil, nil},
			&PackerConfig{
				CorePackerVersionString: lockedVersion,
				Builds: Builds{
					&BuildBlock{
						Sources: []SourceUseBlock{
							{
								SourceRef: SourceRef{
									Type: "null",
									Name: "test",
								},
							},
						},
					},
				},
				Sources: map[SourceRef]SourceBlock{
					{
						Type: "null",
						Name: "test",
					}: {
						Type: "null",
						Name: "test",
					},
				},
				Basedir: filepath.Join("testdata", "datasources"),
				Datasources: Datasources{
					{
						Type: "amazon-ami",
						Name: "test",
					}: {
						Type:   "amazon-ami",
						DSName: "test",
					},
				},
			},
			false, false,
			[]*packer.CoreBuild{
				&packer.CoreBuild{
					Type:           "null.test",
					BuilderType:    "null",
					Builder:        &null.Builder{},
					Provisioners:   []packer.CoreBuildProvisioner{},
					PostProcessors: [][]packer.CoreBuildPostProcessor{},
					Prepared:       true,
					SensitiveVars:  []string{},
				},
			},
			false,
			nil,
		},
		{"recursive datasources",
			defaultParser,
			parseTestArgs{"testdata/datasources/recursive.pkr.hcl", nil, nil},
			&PackerConfig{
				CorePackerVersionString: lockedVersion,
				Builds: Builds{
					&BuildBlock{
						Sources: []SourceUseBlock{
							{
								SourceRef: SourceRef{
									Type: "null",
									Name: "test",
								},
							},
						},
					},
				},
				Sources: map[SourceRef]SourceBlock{
					{
						Type: "null",
						Name: "test",
					}: {
						Type: "null",
						Name: "test",
					},
				},
				Basedir: filepath.Join("testdata", "datasources"),
				Datasources: Datasources{
					{
						Type: "null",
						Name: "foo",
					}: {
						Type:   "null",
						DSName: "foo",
					},
					{
						Type: "null",
						Name: "bar",
					}: {
						Type:   "null",
						DSName: "bar",
					},
					{
						Type: "null",
						Name: "baz",
					}: {
						Type:   "null",
						DSName: "baz",
						Dependencies: []refString{
							{"data", "null", "foo"},
							{"data", "null", "bar"},
						},
					},
					{
						Type: "null",
						Name: "bang",
					}: {
						Type:   "null",
						DSName: "bang",
						Dependencies: []refString{
							{"data", "null", "baz"},
						},
					},
					{
						Type: "null",
						Name: "yummy",
					}: {
						Type:   "null",
						DSName: "yummy",
						Dependencies: []refString{
							{"data", "null", "bang"},
						},
					},
				},
			},
			false, false,
			[]*packer.CoreBuild{
				&packer.CoreBuild{
					Type:           "null.test",
					BuilderType:    "null",
					Builder:        &null.Builder{},
					Provisioners:   []packer.CoreBuildProvisioner{},
					PostProcessors: [][]packer.CoreBuildPostProcessor{},
					Prepared:       true,
					SensitiveVars:  []string{},
				},
			},
			false,
			nil,
		},
		{"untyped datasource",
			defaultParser,
			parseTestArgs{"testdata/datasources/untyped.pkr.hcl", nil, nil},
			&PackerConfig{
				CorePackerVersionString: lockedVersion,
				Basedir:                 filepath.Join("testdata", "datasources"),
			},
			true, true,
			nil,
			false,
			nil,
		},
		{"unnamed source",
			defaultParser,
			parseTestArgs{"testdata/datasources/unnamed.pkr.hcl", nil, nil},
			&PackerConfig{
				CorePackerVersionString: lockedVersion,
				Basedir:                 filepath.Join("testdata", "datasources"),
			},
			true, true,
			nil,
			false,
			nil,
		},
		{"nonexistent source",
			defaultParser,
			parseTestArgs{"testdata/datasources/nonexistent.pkr.hcl", nil, nil},
			&PackerConfig{
				CorePackerVersionString: lockedVersion,
				Basedir:                 filepath.Join("testdata", "datasources"),
				Datasources: Datasources{
					{
						Type: "nonexistent",
						Name: "test",
					}: {
						Type:   "nonexistent",
						DSName: "test",
					},
				},
			},
			true, true,
			nil,
			false,
			nil,
		},
		{"duplicate source",
			defaultParser,
			parseTestArgs{"testdata/datasources/duplicate.pkr.hcl", nil, nil},
			&PackerConfig{
				CorePackerVersionString: lockedVersion,
				Basedir:                 filepath.Join("testdata", "datasources"),
				Datasources: Datasources{
					{
						Type: "amazon-ami",
						Name: "test",
					}: {
						Type:   "amazon-ami",
						DSName: "test",
					},
				},
			},
			true, true,
			nil,
			false,
			nil,
		},
		{"cyclic dependency between data sources",
			defaultParser,
			parseTestArgs{"testdata/datasources/dependency_cycle.pkr.hcl", nil, nil},
			&PackerConfig{
				CorePackerVersionString: lockedVersion,
				Basedir:                 filepath.Join("testdata", "datasources"),
				Datasources: Datasources{
					{
						Type: "null",
						Name: "gummy",
					}: {
						Type:   "null",
						DSName: "gummy",
						Dependencies: []refString{
							{"data", "null", "bear"},
						},
					},
					{
						Type: "null",
						Name: "bear",
					}: {
						Type:   "null",
						DSName: "bear",
						Dependencies: []refString{
							{"data", "null", "gummy"},
						},
					},
				},
			},
			true, true,
			nil,
			false,
			nil,
		},
	}

	testParse(t, tests)
}
