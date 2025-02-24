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
							Documentation: "The position of this hint.\n\nIf multiple hints have the same position, they will be shown in the order\nthey appear in the response.",
						},
						{
							Name: "paddingLeft",
							Type: Type{
								Kind: "base",
								Name: "boolean",
							},
							Optional:      toPtr(true),
							Documentation: "Render padding before the hint.\n\nNote: Padding should use the editor's background color, not the\nbackground color of the hint itself. That means padding can be used\nto visually align/separate an inlay hint.",
						},
						{
							Name: "paddingRight",
							Type: Type{
								Kind: "base",
								Name: "boolean",
							},
							Optional:      toPtr(true),
							Documentation: "Render padding after the hint.\n\nNote: Padding should use the editor's background color, not the\nbackground color of the hint itself. That means padding can be used\nto visually align/separate an inlay hint.",
						},
						{
							Name: "data",
							Type: Type{
								Kind: "reference",
								Name: "LSPAny",
							},
							Optional:      toPtr(true),
							Documentation: "A data entry field that is preserved on an inlay hint between\na `textDocument/inlayHint` and a `inlayHint/resolve` request.",
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
							Optional:      toPtr(true),
							Documentation: "Optional text edits that are performed when accepting this inlay hint.\n\n*Note* that edits are expected to change the document so that the inlay\nhint (or its nearest variant) is now part of the document and the inlay\nhint itself is now obsolete.",
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
							Optional:      toPtr(true),
							Documentation: "The label of this hint. A human readable string or an array of\nInlayHintLabelPart label parts.\n\n*Note* that neither the string nor the label part can be empty.",
						},
						{
							Name: "kind",
							Type: Type{
								Kind: "reference",
								Name: "InlayHintKind",
							},
							Optional:      toPtr(true),
							Documentation: "The kind of this hint. Can be omitted in which case the client\nshould fall back to a reasonable default.",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateStructure(tt.args.s)

			assert.Equal(t, got.String(), tt.want.String())

		})
	}
}

func toPtr(b bool) *bool {
	return &b
}
