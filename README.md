# day-two-concurrency

onion0904の7日間チャレンジの二日目です。
詳しくは[zenn]()

## 作ったもの
resizon(resize+onion)を作成しました。
このプロダクトは与えられたディレクトリの画像のサイズを変更するCLIツールです。
jpegとpngに対応しています。
resizonは一つ一つの画像を並行処理しているため、高速に画像のサイズを変更することができます。


## 実行フロー
```
go run main.go
input create of image directory: ./output //出力先のファイルパス
input path of image directory: ./images //元画像のディレクトリのパス
input size: 100 200//変更するサイズ
```

## サイズについて
### サイズが一つの場合

アスペクト比を維持したまま出力します。
```
input size: 100
```
この時は元画像の短い辺を100pxにして、長い辺を元の画像の比に合わせるように変更します。

### サイズが二つの場合
横、縦の順で入力されます。
```
input size: 100 200
```
横100px,縦200px



## 参考資料
[Goの並行処理入門](https://tech.yappli.io/entry/goroutine-base)

[画像処理](https://blanktar.jp/blog/2024/01/golang-resize-image)