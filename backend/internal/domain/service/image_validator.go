package service

// 画像のバリデーション
type ImageValidator interface {
	// マジックナンバー（バイナリシグネチャ）を使用してファイル形式を検証
	ValidateFormat(data []byte) error

	// ファイルサイズが上限を超えていないかをチェック
	ValidateSize(size int64) error
}
