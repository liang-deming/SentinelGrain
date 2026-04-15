# SentinelGrain

[English](README.md) | [日本語](README.ja.md) | [中文](README.zh.md)

SentinelGrain は、Go と go-zero で実装された軽量なクォータ/レート制御サービスです。
主な機能は次のとおりです。

- クォータ判定を行う RPC サービス（`sentinel-rpc`）
- ルール作成・参照のための Admin HTTP API
- Redis に保存したルールを定期的にメモリへリフレッシュ
- Docker / Kubernetes 向けのデプロイ資産

## ディレクトリ構成

- `rpc/`: gRPC サービス本体、判定ロジック、設定、テスト
- `api/`: クォータルール管理用の Admin REST サービス
- `deploy/`: Docker/Kubernetes 設定と運用メモ
- `common/`: 共通モデル・ユーティリティ

## 前提条件

- Go `1.24+`
- Redis（ローカル開発時の既定は `127.0.0.1:6379`）
- 任意: protobuf 再生成用の `protoc`

## ローカル実行

1. Redis を起動します。
2. RPC サービスを起動します。

```bash
go run ./rpc/cmd/sentinel-rpc -f rpc/etc/sentinel.yaml
```

3. Admin サービスを起動します。

```bash
go run ./api/admin.go -f api/etc/admin.yaml
```

## API / RPC

- Admin API ベース URL: `http://localhost:8888/api/v1`
  - `POST /quota`: ルール作成・更新
  - `GET /quota`: ルール一覧取得
- RPC エンドポイント: `localhost:8080`
  - サービス: `sentinel.Sentinel`
  - メソッド: `Check`

`grpcurl` の例:

```bash
grpcurl -plaintext \
  -d '{"appId":"myapp","resource":"api","dimension":"u1","cost":1}' \
  localhost:8080 sentinel.Sentinel/Check
```

## 設定のポイント

- RPC サーバー設定はトップレベルの `rpc:` 配下にあります。
- Redis 業務設定はトップレベルの `BizRedis:` を使用します。
- `QuotaRefreshInterval` は Redis からの定期リフレッシュ間隔です。
- `L1Cache` はホットパス向けのローカルメモリキャッシュです。

現在の設定例は `rpc/etc/sentinel.yaml` を参照してください。

## デプロイ

Docker/Kubernetes の詳細手順とトラブルシュートは次を参照してください。

- `deploy/README.md`

Docker Desktop + Kubernetes のローカル環境では、ロールアウト時に毎回ユニークタグを使うため、次のスクリプトが便利です。

```bash
./deploy/k8s/deploy-local-image.sh
```

## テスト

```bash
go test ./...
```
