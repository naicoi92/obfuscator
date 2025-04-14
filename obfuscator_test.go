package obfuscator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewObfuscator_Success kiểm tra việc tạo obfuscator thành công
func TestNewObfuscator_Success(t *testing.T) {
	// Arrange & Act
	obf, err := NewObfuscator()

	// Assert
	assert.NoError(t, err, "Không nên có lỗi khi tạo obfuscator")
	assert.NotNil(t, obf, "Obfuscator không nên là nil")
}

// TestNewObfuscator_CachedData kiểm tra việc tạo cached data
func TestNewObfuscator_CachedData(t *testing.T) {
	// Arrange & Act
	obf, err := NewObfuscator()

	// Assert
	require.NoError(t, err, "Không nên có lỗi khi tạo obfuscator")
	assert.NotNil(t, obf.CachedData, "CachedData không nên là nil")
}

// TestNewObfuscator_MultipleObfuscations kiểm tra việc sử dụng obfuscator nhiều lần
func TestNewObfuscator_MultipleObfuscations(t *testing.T) {
	// Arrange
	obf, err := NewObfuscator()
	require.NoError(t, err, "Không nên có lỗi khi tạo obfuscator")

	// Mã JavaScript đơn giản để test
	jsCode := "function test() { return 'Hello World'; }"

	// Act & Assert
	// Thực hiện obfuscate nhiều lần
	for i := 0; i < 3; i++ {
		result, err := obf.Obfuscate(jsCode)
		assert.NoError(t, err, "Obfuscation #%d không nên có lỗi", i+1)
		assert.NotEmpty(t, result, "Kết quả obfuscation #%d không nên rỗng", i+1)
	}
}
