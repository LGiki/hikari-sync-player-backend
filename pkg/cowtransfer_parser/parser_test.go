package cowtransfer_parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCowTransferFileId(t *testing.T) {
	fileId, err := getCowTransferFileId("https://cowtransfer.com/s/example")
	assert.Nil(t, err)
	assert.Equal(t, "example", fileId)
	fileId, err = getCowTransferFileId("http://cowtransfer.com/s/example")
	assert.Nil(t, err)
	assert.Equal(t, "example", fileId)
	fileId, err = getCowTransferFileId("https://cowtransfer.com/s/example/")
	assert.Nil(t, err)
	assert.Equal(t, "example", fileId)
	fileId, err = getCowTransferFileId("http://cowtransfer.com/s/example/")
	assert.Nil(t, err)
	assert.Equal(t, "example", fileId)
	fileId, err = getCowTransferFileId("https://cowtransfer.com/example/")
	assert.NotNil(t, err)
	assert.Equal(t, "", "")
}
