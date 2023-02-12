# auth-test-go

* goでの認証処理を勉強するプロジェクト
* 認証方式として下記の2つを実装
  1. DBにトークンを保存してアクセスの都度検証する方式
     * トークンに有効期限を設定して、切れたら401を返す
  2. JWTによりIDトークンとリフレッシュトークンを発行してIDトークンに包含される有効期限を検証して認証する方式
     * リフレッシュトークンは1と同じでDBに登録する
     * IDトークンの有効期限が切れた時のみリフレッシュトークンを用いて新しいIDトークンを発行
     * OpenID connectにおけるID Providerに相当する処理を実装

## Swagger

* 実行前にswaggerを生成してください
  * Dockerfileには組み込み済み
* gin-swagger生成における注意点は以下の通りです
  * 1つのハンドラ関数を2つのエンドポイントで使用してます
  * ただし、swaggerの仕様で1つのハンドラ関数に定義できるエンドポイントのswaggerコメントは1つまでです
  * そのため試行するエンドポイント毎に書き直し、swaggerを再生成する必要があります
  * ./infra/controller/user_accounts.goの`@Router`を試行するエンドポイントに合わせてください
  * `@Router`を変更するのは以下の関数です
    * Get
    * Update
    * Delete

```
# session or auth のいずれかでswaggerを生成します
# /v1/session/users を試行する場合は下記のように記述
# @Router /session/users/

# /v1/auth/session を試行する場合は下記のように記述
# @Router /auth/users/

# ローカルでswaggerを生成する場合
$ go get -u github.com/swaggo/swag/cmd/swag
$ go install github.com/swaggo/swag/cmd/swag@v1.8.0
$ swag init
```

## 実行方法

```
# 環境変数設定ファイル用意
$ vi environment.txt
# MySQLのrootパスワード
PASSWORD=
ENCRYPT_SECRET=
$ export $(cat environment.txt | grep -v ^#)

# docker-compose.yamlには下記のイメージ名でAPIを起動
$ docker build -t auth-test-go .
$ docker compose up

# MySQLへのDB初期化
# compose後に別ターミナルを使用する場合、再度環境変数を読み込む
$ export $(cat environment.txt | grep -v ^#)
$ cd migration
# マイグレーション実行
$ go run main.go
```

## アクセス方法

- ブラウザで下記のURLでswagger UIにアクセス
  - http://localhost:8080/v1/swagger/index.html
- `/v1/session/users/`の場合は Login/Logout を使用してください
- `/v1/auth/users/`の場合は Claim/Refresh を使用してください

## 注意点

1. リクエストボディのフォーマットに全角文字が存在する場合にpanicを起こす問題が未解決
