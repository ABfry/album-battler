// scanner.go は repository 実装から共有される rowScanner インターフェースを定義する。
// 目的: *sql.Row と *sql.Rows の双方を一貫したシグネチャで扱い、エンティティ変換ロジックを
// 1 箇所に閉じ込めること。
package mysql

// rowScanner は database/sql が返す *sql.Row / *sql.Rows をまとめて扱うための最小インターフェース。
// why: createXxx 系の関数で同じ Scan 呼び出しを共有したいが、型が異なると共通化できないため。
type rowScanner interface {
	Scan(dest ...any) error
}
