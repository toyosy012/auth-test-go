# auth-test-go

* goでの認証処理を勉強するプロジェクト
* 認証方式として下記の2つを実装
  1. DBにトークンを保存してアクセスの都度検証する方式
     * トークンに有効期限を設定して、切れたら401を返す
  2. JWTによりIDトークンとリフレッシュトークンを発行してIDトークンに包含される有効期限を検証して認証する方式
     * リフレッシュトークンは1と同じでDBに登録する
     * IDトークンの有効期限が切れた時のみリフレッシュトークンを用いて新しいIDトークンを発行
     * OpenID connectにおけるID Providerに相当する処理を実装

## 実行方法

```
# 環境変数設定ファイル用意
$ vi environment.txt
PASSWORD=        // MySQLのパスワード
ENCRYPT_SECRET=
EMAIL=
USER_PASSWORD=  // Migration時のユーザアカウントのパスワード
USER_NAME=
$ export $(cat environment.txt | grep -v ^#)

# docker-compose.yamlで下記のイメージ名で起動
$ docker build -t auth-test-go .
$ docker compose up

# MySQLへのDB初期化
# docker composeにより別のターミナルを使用する場合には再度環境変数を読み込む
$ $ export $(cat environment.txt | grep -v ^#)
$ cd migration
$ go run main.go

$ curl http://0.0.0.0:8080/v1/users
```

## API

### API一覧

- `password`には乱数生成した文字列を使用してください

```
GET    /v1/users
POST   /v1/users/new
-H "content-type: application/json" -d '{"email": "test@example.com", "password": "", "name": "test"}'

# セッション認証
# トークンはレスポンスの`Authorization`ヘッダーから取得
POST   /v1/session/login
-v -H "content-type: application/json" -d '{"email": "test@example.com", "password": ""}'
DELETE /v1/session/logout/:id
-H "Authorization: Bearer `session token`"
GET    /v1/session/users/:id
-H "Authorization: Bearer `session token`"
PUT    /v1/session/users/:id
-H "Authorization: Bearer `session token`" -H "content-type: application/json" -d '{"email": "test@example.com", "password": "", "name": "test_updated"}'
DELETE /v1/session/users/:id
-H "Authorization: Bearer `session token`"

# oauth認証
POST   /v1/oauth/claim
-H "content-type: application/json" -d '{"email": "test@example.com", "password": ""}'
# claimから得たリフレッシュトークンをvalueに渡す(デフォルト1h有効)
POST   /v1/oauth/refresh
-H "content-type: application/json" -d '{"value": ""}'
GET    /v1/oauth/users/:id
-H "Authorization: Bearer `access oauth token`"
PUT    /v1/oauth/users/:id
-H "content-type: application/json" -H "Authorization: Bearer `access oauth token`" -d '{"email": "test@example.com", "password": "", "name": "test_updated"}'
DELETE /v1/oauth/users/:id
-H "Authorization: Bearer `access oauth token`"
```

## 注意点

1. リクエストボディのフォーマットに全角文字が存在する場合にpanicを起こす問題が未解決

## 
