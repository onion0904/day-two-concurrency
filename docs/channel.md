# channel
channelはgoroutine間の通信に利用されるもの。

# channelとsyncの違い

## sync
共有メモリやグローバル変数を使い、データにアクセス。
低水準の同期手段であり、正しく使わないとバグを生みやすい。(sync.Onceとsync.WaitGroupは高水準)
高パフォーマンスが求められるような場面では制御が明確にできるため、利用される。

## channel
通信によってメモリを共有して、データにアクセス。
Goでは「共有メモリによって通信するのではなく、通信によってメモリを共有せよ」という言葉があり、こちらが推奨されてる。
複数のgoroutineを協調動作させたり、所有権の変更したりする場面で利用される。

# 双方向チャネル
送受信可能
以下のように宣言
```
// interface{}型のchannelを宣言
var dataStream chan interface{}
// 初期化
dataStream = make(chan interface{})
```
この場合は任意の値を書き込み読み込みができる。

# 単方向チャネル
受信用と送信用の二つがあります。
受信用は<-chanと書く
```
var dataStream <-chan interface{}
dataStream = make(<-chan interface{})
```
送信用はchan<-と書く
```
var dataStream chan<- interface{}
dataStream = make(chan<- interface{})
```
ちなみに、双方向チャネルに代入することもできる。

使い方は以下のようにする
```
package main

import "fmt"

func main() {
    stringStream := make(chan string)

    go func() {
        stringStream <- "Hello channel" //1
    }()

    salutation, ok := <-stringStream //2
    fmt.Printf("(%v): %v", ok, salutation)
}
// 結果: (true): Hello channel
```
1. 書き込み
2. valueとbool(channelが閉じられたかどうか)を返す。

明示的に閉じるには以下のように書く
```
defer close(intStream)
```

- 閉じたとしても、受信できる。
- closeすることで親のブロックを解除できる(channelは自動で書き込みされるまで親をブロックする)
- 複数の送受信が可能でvalue := range chで取得できる。同じ変数の受信者が複数あるときはランダムになる

# バッファ付きチャネル
以下のように書く
```
intStream := make(chan int, 4)
```
- 明示的にいくつ送受信できるかが分かる
- バッファが満杯なるまでブロックせずに送受信できる
- 満杯まで書き込まれたらすぐに読み込みがされる
- 大きすぎても小さすぎても良くない(メモリを用意するので)

# 注意点
- 書き込むゴルーチンとチャネルを読む込むゴルーチンとでしっかりと責務を分けて実装すること
- 書き込み側
    - チャネルの初期化
    - 書き込みを行う、もしくは他のゴルーチンへ所有権を渡す
    - チャネルを閉じる
    - 上記の3つをカプセル化する
-  読み込み側
    - チャネルがいつ閉じられたか
    - ブロックする操作は慎重に扱う

# ブロックの基本ルール
- 非バッファチャネル（make(chan T)）
    - 送信側（ch <- val）は、受信者がいなければブロック。
    - 受信側（val := <-ch）は、送信者がいなければブロック。

- バッファチャネル（make(chan T, N)）
    - 送信側は、バッファが満杯になるまでブロックしない。
    - 受信側は、バッファに何もなければブロック。