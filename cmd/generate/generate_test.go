package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateStructure(t *testing.T) {
	type args struct {
		s Structure
	}
	tests := []struct {
		name string
		args args
		want *bytes.Buffer
	}{
		{
			name: "Verify structure with array",
			args: args{
				s: Structure{
					Name: "InlayHint",
					Properties: []Properties{
						{
							Name: "position",
							Type: Type{
								Kind: "reference",
								Name: "Position",
							},
							Documentation: "The position of this hint.\n\n" +
								"If multiple hints have the same position, they will be shown in the order\n" +
								"they appear in the response.",
						},
						{
							Name: "paddingLeft",
							Type: Type{
								Kind: "base",
								Name: "boolean",
							},
							Optional: toPtr(true),
							Documentation: "Render padding before the hint.\n\n" +
								"Note: Padding should use the editor's background color, not the\n" +
								"background color of the hint itself. That means padding can be used\n" +
								"to visually align/separate an inlay hint.",
						},
						{
							Name: "paddingRight",
							Type: Type{
								Kind: "base",
								Name: "boolean",
							},
							Optional: toPtr(true),
							Documentation: "Render padding after the hint.\n\n" +
								"Note: Padding should use the editor's background color, not the\n" +
								"background color of the hint itself. That means padding can be used\n" +
								"to visually align/separate an inlay hint.",
						},
						{
							Name: "data",
							Type: Type{
								Kind: "reference",
								Name: "LSPAny",
							},
							Optional: toPtr(true),
							Documentation: "A data entry field that is preserved on an inlay hint between\n" +
								"a `textDocument/inlayHint` and a `inlayHint/resolve` request.",
						},
						{
							Name: "textEdits",
							Type: Type{
								Kind: "array",
								Element: &Option{
									Kind: "reference",
									Name: "TextEdit",
								},
							},
							Optional: toPtr(true),
							Documentation: "Optional text edits that are performed when accepting this inlay hint.\n\n" +
								"*Note* that edits are expected to change the document so that the inlay\n" +
								"hint (or its nearest variant) is now part of the document and the inlay\n" +
								"hint itself is now obsolete.",
						},
						{
							Name: "tooltip",
							Type: Type{
								Kind: "or",
								Items: []Option{
									{
										Kind: "base",
										Name: "string",
									},
									{
										Kind: "reference",
										Name: "MarkupContent",
									},
								},
							},
							Optional:      toPtr(true),
							Documentation: "The tooltip text when you hover over this item.",
						},
						{
							Name: "label",
							Type: Type{
								Kind: "or",
								Items: []Option{
									{
										Kind: "base",
										Name: "string",
									},
									{
										Kind: "array",
										// Element: Element{
										// 	Kind: "reference",
										// 	Name: "InlayHintLabelPart",
										// },
									},
								},
							},
							Optional: toPtr(true),
							Documentation: "The label of this hint. A human readable string or an array of\n" +
								"InlayHintLabelPart label parts.\n\n" +
								"*Note* that neither the string nor the label part can be empty.",
						},
						{
							Name: "kind",
							Type: Type{
								Kind: "reference",
								Name: "InlayHintKind",
							},
							Optional: toPtr(true),
							Documentation: "The kind of this hint. Can be omitted in which case the client\n" +
								"should fall back to a reasonable default.",
						},
					},
					// Documentation: "Inlay hint information.\n\n@since 3.17.0",
					// Since:         "3.17.0",
				},
			},
			want: func() *bytes.Buffer {
				w, err := os.ReadFile("testdata/InlayHintStructure.txt")
				if err != nil {
					t.Fatal(err)
				}
				return bytes.NewBuffer(w)
			}(),
		},
		{
			name: "Verify structure with mixins",
			args: args{
				s: Structure{
					Name: "SemanticTokensRangeParams",
					Properties: []Properties{
						{
							Name: "textDocument",
							Type: Type{
								Kind: "reference",
								Name: "TextDocumentIdentifier",
							},
							Documentation: "The text document.",
						},
						{
							Name: "range",
							Type: Type{
								Kind: "reference",
								Name: "Range",
							},
							Documentation: "The range the semantic tokens are requested for.",
						},
					},
					Mixins: []Option{
						{
							Kind: "reference",
							Name: "WorkDoneProgressParams",
						},
						{
							Kind: "reference",
							Name: "PartialResultParams",
						},
					},

					// Documentation: "@since 3.16.0",
					// Since:         "3.16.0",
				},
			},
			want: func() *bytes.Buffer {
				w, err := os.ReadFile("testdata/SemanticTokensRangeParams.txt")
				if err != nil {
					t.Fatal(err)
				}
				return bytes.NewBuffer(w)
			}(),
		},
		{
			name: "Verify structure with mixins and extends",
			args: args{
				s: Structure{
					Name: "LinkedEditingRangeParams",
					Mixins: []Option{
						{
							Kind: "reference",
							Name: "WorkDoneProgressParams",
						},
					},
					Extends: []Option{
						{
							Kind: "reference",
							Name: "TextDocumentPositionParams",
						},
					},
					// Documentation: "@since 3.16.0",
					// Since:         "3.16.0",
				},
			},
			want: func() *bytes.Buffer {
				w, err := os.ReadFile("testdata/LinkedEditingRangeParams.txt")
				if err != nil {
					t.Fatal(err)
				}
				return bytes.NewBuffer(w)
			}(),
		},
		{
			name: "Verify unexported structure is skipped",
			args: args{
				s: Structure{
					Name:       "_UnExported",
					Properties: []Properties{},
				},
			},
			want: nil,
		},
		{
			name: "Verify ToTitleCase handles UTF-8",
			args: args{
				s: Structure{
					Name: "Utf8Struct",
					Properties: []Properties{
						{
							Name: "café",
							Type: Type{
								Kind: "base",
								Name: "string",
							},
						},
					},
				},
			},
			want: bytes.NewBufferString("type Utf8Struct struct {\n\tCafé string\n}\n"),
		},
		{
			name: "Verify optional non-array becomes pointer",
			args: args{
				s: Structure{
					Name: "Optional",
					Properties: []Properties{
						{
							Name:     "count",
							Optional: toPtr(true),
							Type: Type{
								Kind: "base",
								Name: "integer",
							},
						},
					},
				},
			},
			want: bytes.NewBufferString("type Optional struct {\n\tCount *int\n}\n"),
		},
		{
			name: "Verify optional array is not pointer",
			args: args{
				s: Structure{
					Name: "OptionalArray",
					Properties: []Properties{
						{
							Name:     "items",
							Optional: toPtr(true),
							Type: Type{
								Kind:    "array",
								Element: &Option{Kind: "base", Name: "string"},
							},
						},
					},
				},
			},
			want: bytes.NewBufferString("type OptionalArray struct {\n\tItems []string\n}\n"),
		},
		{
			name: "Verify hidden mixins are skipped",
			args: args{
				s: Structure{
					Name: "WithHiddenMixins",
					Mixins: []Option{
						{Kind: "reference", Name: "_Hidden"},
						{Kind: "reference", Name: "Visible"},
					},
				},
			},
			want: bytes.NewBufferString("type WithHiddenMixins struct {\n\tVisible\n}\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateStructure(tt.args.s)
			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			assert.NotNil(t, got)
			assert.Equal(t, tt.want.String(), got.String())
		})
	}
}

func TestConvertType(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "boolean", in: "boolean", want: "bool"},
		{name: "uinteger", in: "uinteger", want: "uint"},
		{name: "integer", in: "integer", want: "int"},
		{name: "decimal", in: "decimal", want: "float64"},
		{name: "LSPAny", in: "LSPAny", want: "interface{}"},
		{name: "URI", in: "URI", want: "string"},
		{name: "DocumentUri", in: "DocumentUri", want: "string"},
		{name: "empty", in: "", want: "interface{}"},
		{name: "passthrough", in: "CustomType", want: "CustomType"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ConvertType(tt.in))
		})
	}
}

func toPtr(b bool) *bool {
	return &b
}
