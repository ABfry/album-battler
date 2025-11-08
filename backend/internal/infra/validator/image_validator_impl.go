package validator

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

// 画像のバリデーション
type imageValidator struct {
	maxSize        int64
	allowedFormats []string
}

var (
	webpStartBytes = []byte{0x52, 0x49, 0x46, 0x46}                                                 // 'R' 'I' 'F' 'F'
	webpFormatChar = []byte{0x57, 0x45, 0x42, 0x50}                                                 // 'W' 'E' 'B' 'P'
	heifHeaderChar = []byte{0x6d, 0x69, 0x66, 0x31, 0x68, 0x65, 0x69, 0x63, 0x68, 0x65, 0x76, 0x63} // "mif1heichevc"
)

// maxSize: 許可する最大ファイルサイズ（バイト単位）
func NewImageValidator(maxSize int64) service.ImageValidator {
	return &imageValidator{
		maxSize: maxSize,
		allowedFormats: []string{
			"image/jpeg",
			"image/png",
			"image/webp",
			"image/heic",
		},
	}
}

// マジックナンバーを使用して画像形式を検証する
// サポート形式: JPEG, PNG, WebP, HEIC
func (v *imageValidator) ValidateFormat(data []byte) error {
	if len(data) < 12 {
		return errors.New("invalid image: file too small")
	}

	if isJpeg(data) {
		return nil
	}

	if isPng(data) {
		return nil
	}

	if isWebp(data) {
		return nil
	}

	if isHeif(data) {
		return nil
	}

	return errors.New("unsupported image format: only JPEG, PNG, WebP, and HEIC are allowed")
}

// JPEGフォーマットかどうかを判定する
func isJpeg(imageBytes []byte) bool {
	return bytes.HasPrefix(imageBytes, []byte{0xFF, 0xD8, 0xFF})
}

// PNGフォーマットかどうかを判定する
func isPng(imageBytes []byte) bool {
	return bytes.HasPrefix(imageBytes, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
}

// HEIFフォーマットかどうかを判定する
func isHeif(imageBytes []byte) bool {
	// mif1heichevc in 17-28 byte
	if len(imageBytes) < 28 {
		return false
	}
	if simpleByteRangeEqual(imageBytes[16:28], heifHeaderChar, 12) {
		return true
	}
	return false
}

// WebPフォーマットかどうかを判定する
func isWebp(imageBytes []byte) bool {
	if len(imageBytes) < 12 {
		return false
	}
	// 最初の4バイトをチェック
	if !simpleByteRangeEqual(imageBytes, webpStartBytes, len(webpStartBytes)) {
		return false
	}
	// 9から12バイトをチェック
	if !simpleByteRangeEqual(imageBytes[8:12], webpFormatChar, 4) {
		return false
	}
	return true
}

// バイト配列の範囲比較
func simpleByteRangeEqual(bytes1, bytes2 []byte, checkLength int) bool {
	if len(bytes1) < checkLength || len(bytes2) < checkLength {
		return false
	}
	for index := 0; index < checkLength; index++ {
		if bytes1[index] != bytes2[index] {
			return false
		}
	}
	return true
}

// ファイルサイズが上限を超えていないかを検証する
func (v *imageValidator) ValidateSize(size int64) error {
	if size > v.maxSize {
		return fmt.Errorf("image size (%d bytes) exceeds maximum allowed size (%d bytes)", size, v.maxSize)
	}
	if size == 0 {
		return errors.New("image size is zero")
	}
	return nil
}
