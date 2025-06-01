# select
switch文のようなものであり、channelの送受信の管理をするハブのようなもの。
複数のチャネル間での通信を統合し、キャンセル処理、タイムアウト、待機、デフォルト処理などをまとめて行うことができる。
以下のように書く
```
func main() {
    var c1,c2<-chan any
    var c3 chan<-any

    select {
    case <-c1:
    // do something
    case <-c2:
    // do something
    case c3<-struct{}{}:
    // do something
    }
}
```

# タイムアウト
一つも読み込めない場合、永遠にブロックされるのでタイムアウトしたい。
以下のように書く。
```
func main() {
    var c1, c2 <-chan int

    select {
    case <-c1:
    case <-c2:
    case <-time.After(2 * time.Second):
        fmt.Println("timeout")
    }
}
```
time.After関数はtime.Durationを引数に取り、引数で与えた経過時間後に現在時刻を送信するチャネルを返すので、どのチャネルも読み込めなくてタイムアウトさせたいときに便利です。

# デフォルト節
読み込まれない間にタイムアウトではなく処理をさせたい場合、
以下のように書き、インクリメントさせたりなどできる。
```
select {
        case <-channel1:
            fmt.Println("channel1を受信")
            break loop
        default:
            count++
            time.Sleep(1 * time.Second)
        }
```