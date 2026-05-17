# hateb

指定したはてな ID の公開はてなブックマーク RSS を取得し、ブックマークタイトルと URL を出力する CLI ツールです。

## 使い方

```sh
go run . --user <はてなID>
go run . --user <はてなID> --since <YYYY-MM-DD>
```

ビルドする場合:

```sh
go build -o hateb
./hateb --user <はてなID>
./hateb --user <はてなID> --since <YYYY-MM-DD>
```

`--since` を指定した場合は、指定日以降のブックマークだけを出力します。
未指定の場合は、RSS で取得できる 1 ページ目の内容を出力します。

出力形式は 1 件につき日付とタイトル、URL、空行です。
タイトルはターミナル幅に合わせ、入りきらない場合は末尾を `...` で省略します。

```text
yyyy/MM/dd title
           url
```
