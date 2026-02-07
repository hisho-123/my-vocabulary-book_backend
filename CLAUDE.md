# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

単語帳アプリケーションのバックエンドAPI。Goで構築され、アプリの規模に適した簡略化されたクリーンアーキテクチャパターンを採用。

## ファイル構造

- リポジトリルート: 開発ツール（docker-compose.yml等）を配置
- src/: アプリケーションコード（Goソース、インフラ実装）を配置

## 技術スタック

- 言語: Go 1.23.3
- フレームワーク: Gin (HTTPルーター)
- データベース: MySQL
- 認証: JWT トークン (golang-jwt/jwt)
- パスワードハッシュ化: bcrypt

## 開発コマンド

### アプリケーションの起動

```bash
cd src
go run main.go
```

サーバーは `localhost:8080` で起動、ベースパスは `/api`。

### データベースセットアップ

```bash
# MySQLコンテナを起動（リポジトリルートから実行）
docker compose up -d

# データベースに直接アクセス
docker exec -it db bash
mysql -u root -p  # password: root

# または直接アクセス
docker exec -it db mysql -uroot -proot

# MySQL内での操作
show databases;
use my_vocabulary_book;
show tables;
```

データベース接続: `root:root@tcp(127.17.0.1:3306)/my_vocabulary_book`

スキーマファイル: `src/infra/db/migrations/schema.sql`

### ビルドとテスト

```bash
cd src
go build -o src
go test ./...
```

### モジュール管理

```bash
cd src
go mod tidy
go mod download
```

## アーキテクチャ

簡略化されたクリーンアーキテクチャを採用。4つの主要レイヤーで構成:

### 1. Domain層 (`src/domain/`)
- コアエンティティとモデル (`model.go`)
- 認証ユーティリティ (`authUser.go`): パスワードハッシュ化、JWTトークンの生成/検証
- ステータスコード定数 (`statusCode.go`): HTTPステータスコードの集中管理
- 外部依存なし（純粋なビジネスロジック）

### 2. Usecase層 (`src/usecase/`)
- ビジネスロジックのオーケストレーション
- domainのユーティリティとgatewayのメソッドを呼び出す
- ファイル: `user.go`, `create.go`, `read.go`, `delete.go`

### 3. Interface層 (`src/interface/`)
- **Controller** (`interface/controller/`): HTTPリクエストハンドラー、バリデーション、レスポンス整形
- **Gateway** (`interface/gateway/`): データベース操作とデータ変換
  - 注意: 従来のクリーンアーキテクチャと異なり、gatewayはデータ取得とデータ整形の両方を行う（コードの複雑さを避けるため）

### 4. Infrastructure層 (`src/infra/`)
- 外部システムとの統合
- `router.go`: Ginルーターのセットアップ、CORS設定、ルート定義
- `db/db.go`: データベース接続管理
- `db/migrations/`: SQLスキーマファイル

### データフロー

リクエスト → Controller (バリデーション) → Usecase (ビジネスロジック) → Gateway (DB操作) → レスポンス

## エラーハンドリングパターン

エラーは `domain/statusCode.go` の文字列定数を使って伝播:

```go
// エラーを文字列として返す
return fmt.Errorf(domain.Unauthorized)
return fmt.Errorf(domain.InternalServerError)
```

ControllerはこれらをHTTPステータスコードにマッピング (`controller/statusCode.go`):

```go
func statusCode(err error) int {
    switch err.Error() {
    case domain.BadRequest:
        return http.StatusBadRequest
    // ... など
    }
}
```

主なステータスコード:
- 400 (BadRequest): 入力不正、JSON解析エラー
- 401 (Unauthorized): トークンなし/無効、認証失敗
- 403 (Forbidden): 認可失敗（他ユーザーのリソースへのアクセス）
- 404 (NotFound): リソースが見つからない
- 409 (Conflict): ユーザー名の重複
- 422 (UnprocessableEntity): バリデーション失敗（文字数制限）
- 500 (InternalServerError): サーバー/DBエラー

## 認証フロー

1. ユーザー登録/ログイン時にJWTトークンとuserIdを返す
2. 認証が必要なリクエストでは `Token` ヘッダーでトークンを渡す
3. トークン検証時にclaimsからuserIdを抽出
4. Controllerでトークンチェック、Usecaseでclaimsを使って認可処理

JWTトークンの有効期限は24時間。秘密鍵は環境変数 `JWT_KEY` に保存。

## データベーススキーマ

カスケード削除設定の3つのテーブル:
- `users`: user_id (PK), user_name (unique, 最大50文字), password (ハッシュ化, 最大70文字)
- `books`: book_id (PK), user_id (FK), book_name (最大20文字), first_review
- `words`: word_id (PK), book_id (FK), word (最大50文字), translated_word (最大50文字)

## APIエンドポイント

全エンドポイントは `/api` プレフィックス付き:

**認証:**
- POST `/register` - ユーザー登録
- POST `/login` - ユーザー認証
- GET `/home` - トークン検証チェック

**ユーザー管理:**
- DELETE `/user-delete` - ユーザー削除（トークン必須）

**単語帳:**
- POST `/book` - 単語帳と単語を作成
- GET `/book?bookId=X` - 単語帳の詳細と単語を取得
- GET `/book-list` - 認証ユーザーの全単語帳を取得
- DELETE `/book-delete` - 単語帳を削除

## 重要な注意事項

- CORSは `http://localhost:5173` （フロントエンド開発サーバー）用に設定 - デプロイ時に変更必須
- データベース接続はリクエストごとにオープン/クローズ（現在コネクションプーリングなし）
- Gatewayがデータ変換を実行（簡潔さのため厳密なクリーンアーキテクチャから逸脱）
- セキュリティ: ログイン時にユーザーが存在しない場合も401を返す（パスワード間違いと同じ）- ユーザー名列挙攻撃を防ぐため
