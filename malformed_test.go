package asyncapi_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/psyb0t/ctxerrors/commerr"
	asyncapi "github.com/psyb0t/go-asyncapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const documentHeader = "asyncapi: 3.0.0\ninfo: {title: t, version: '1'}\n"

func writeDocument(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "asyncapi.yaml")
	require.NoError(t, os.WriteFile(path, []byte(documentHeader+body), 0o600))

	return path
}

func TestLoadRejectsMalformedSections(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		body    string
		wantErr error
	}{
		{
			name:    "server is not a mapping",
			body:    "servers: {s: 1}\n",
			wantErr: commerr.ErrInvalidValue,
		},
		{
			name:    "server ref points nowhere",
			body:    "servers: {s: {$ref: '#/nope'}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "server variable is not a mapping",
			body:    "servers: {s: {host: h, protocol: ws, variables: {v: 1}}}\n",
			wantErr: commerr.ErrInvalidValue,
		},
		{
			name:    "server variable ref points nowhere",
			body:    "servers: {s: {host: h, protocol: ws, variables: {v: {$ref: '#/nope'}}}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "channel is not a mapping",
			body:    "channels: {c: 1}\n",
			wantErr: commerr.ErrInvalidValue,
		},
		{
			name:    "channel ref points nowhere",
			body:    "channels: {c: {$ref: '#/nope'}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "parameter is not a mapping",
			body:    "channels: {c: {address: 'a/{p}', parameters: {p: 1}}}\n",
			wantErr: commerr.ErrInvalidValue,
		},
		{
			name:    "parameter ref points nowhere",
			body:    "channels: {c: {address: 'a/{p}', parameters: {p: {$ref: '#/nope'}}}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "message is not a mapping",
			body:    "channels: {c: {messages: {m: 1}}}\n",
			wantErr: commerr.ErrInvalidValue,
		},
		{
			name:    "message ref points nowhere",
			body:    "channels: {c: {messages: {m: {$ref: '#/nope'}}}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "payload ref points nowhere",
			body:    "channels: {c: {messages: {m: {payload: {$ref: '#/nope'}}}}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "headers ref points nowhere",
			body:    "channels: {c: {messages: {m: {headers: {$ref: '#/nope'}}}}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "correlationId is not a mapping",
			body:    "channels: {c: {messages: {m: {correlationId: 1}}}}\n",
			wantErr: commerr.ErrInvalidValue,
		},
		{
			name:    "correlationId ref points nowhere",
			body:    "channels: {c: {messages: {m: {correlationId: {$ref: '#/nope'}}}}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "message trait ref points nowhere",
			body:    "channels: {c: {messages: {m: {traits: [{$ref: '#/nope'}]}}}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "component message is not a mapping",
			body:    "components: {messages: {m: 1}}\n",
			wantErr: commerr.ErrInvalidValue,
		},
		{
			name:    "operation is not a mapping",
			body:    "operations: {o: 1}\n",
			wantErr: commerr.ErrInvalidValue,
		},
		{
			name:    "operation ref points nowhere",
			body:    "operations: {o: {$ref: '#/nope'}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "operation channel ref points nowhere",
			body:    "operations: {o: {action: send, channel: {$ref: '#/channels/nope'}}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "operation trait ref points nowhere",
			body:    "operations: {o: {action: send, traits: [{$ref: '#/nope'}]}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "operation message ref points nowhere",
			body:    "channels: {c: {address: c}}\noperations: {o: {action: send, channel: {$ref: '#/channels/c'}, messages: [{$ref: '#/nope'}]}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "reply ref points nowhere",
			body:    "channels: {c: {address: c}}\noperations: {o: {action: send, channel: {$ref: '#/channels/c'}, reply: {$ref: '#/nope'}}}\n",
			wantErr: commerr.ErrNotFound,
		},
		{
			name:    "reply channel ref points nowhere",
			body:    "channels: {c: {address: c}}\noperations: {o: {action: send, channel: {$ref: '#/channels/c'}, reply: {channel: {$ref: '#/channels/nope'}}}}\n",
			wantErr: commerr.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := asyncapi.Load(writeDocument(t, tc.body))
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestLoadRejectsRootThatIsNotAMapping(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "list.yaml")
	require.NoError(t, os.WriteFile(path, []byte("- a\n"), 0o600))

	_, err := asyncapi.Load(path)
	require.ErrorIs(t, err, commerr.ErrParseFailed)
}

func TestLoadMinimalDocument(t *testing.T) {
	t.Parallel()

	document, err := asyncapi.Load(writeDocument(t, ""))
	require.NoError(t, err)

	assert.Empty(t, document.Servers)
	assert.Empty(t, document.Channels)
	assert.Empty(t, document.Operations)
	assert.Empty(t, document.Messages)
}

func TestLoadKeepsBooleanPayloadSchema(t *testing.T) {
	t.Parallel()

	document, err := asyncapi.Load(
		writeDocument(t, "channels: {c: {messages: {m: {payload: true}}}}\n"),
	)
	require.NoError(t, err)
	require.Len(t, document.Messages, 1)

	assert.Equal(t, true, document.Messages[0].Payload)
}

func TestLocationLastTokenOfRoot(t *testing.T) {
	t.Parallel()

	assert.Empty(t, asyncapi.Location{File: "x.yaml"}.LastToken())
}
