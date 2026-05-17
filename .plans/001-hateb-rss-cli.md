# hateb RSS CLI 実装計画

## 目的

指定したはてな ID の公開はてなブックマーク RSS を取得し、ブックマークタイトルと URL を標準出力へ出力する CLI ツールを実装する。`--from` 指定時は指定日以降に絞り込み、未指定時は RSS で取得できる 1 ページ目の内容を出力する。

## コマンド仕様

```sh
hateb --user <はてなID> --from <yyyyMMdd>
```

`--from` は任意指定とする。

## 修正方針

1. ルートの `main.go` に CLI エントリポイントを置き、起動処理だけにする。
2. `internal/cli` で `flag` パッケージによる `--user` と `--from` の解釈、入出力、終了コード制御を扱う。
3. `internal/hateb` で RSS URL 生成、HTTP 取得、RSS item の XML 解析を扱う。
4. `--from` が指定された場合は `yyyyMMdd` として検証し、日付として解釈する。
5. RSS の取得先は `https://b.hatena.ne.jp/{user}/rss` とする。
6. 標準ライブラリの XML デコードで RSS item を読み取る。
7. `--from` が指定された場合のみ、item の日付を解析して指定日以降のものだけを出力する。
8. `--from` が未指定の場合は、RSS で取得した 1 ページ目の item をそのまま出力する。
9. 出力は 1 件ごとに `yyyy/MM/dd title`、次行に 11 桁インデントした URL、空行の形式にする。
10. タイトルはターミナル幅に合わせ、1 行目がはみ出す場合は末尾を `...` で省略する。
11. 入力不備、HTTP エラー、RSS 解析エラーは標準エラーへ出し、非 0 終了にする。
12. 日付判定や RSS 解析の主要処理にテストを追加する。

## 想定編集ファイル

- `main.go`
- `internal/cli/*.go`
- `internal/cli/*_test.go`
- `internal/hateb/*.go`
- `internal/hateb/*_test.go`
- `README.md`

## 確認方法

- `go test ./...`
- ローカルで `go run . --user <はてなID> --from <yyyyMMdd>` を実行し、取得と出力を確認する。
