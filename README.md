# auth-test-go
goでの認証処理を勉強するプロジェクト

## DB

### docker compose

```
$ docker compose up
$ docker compose down
```

### 自動マイグレーション

- ./migration/main.goを用いて自動マイグレーションを行う

```
$ export $(cat environment.txt | grep -v ^#)
$ cd migration
$ go run main.go
```
