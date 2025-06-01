# コルーチン
プログラマ自身で実行中に一時停止し、任意の箇所で再開することができる。

# ゴルーチン
OSによってタスクやプロセスが強制中断されず、タスク地震が自発的に制御を放棄するまで実行を続ける(非プリエンプティブ)。割り込み処理をされることが無い。
ゴルーチンがブロックしたら自動的に一時停止し、ブロックが開放されたら再開するようにGoのランタイムよってのみ制御される。
非同期的な実行方法は以下のように、goの後に関数を呼び出すだけ
```
go hello()
```

# Goランタイムとは
Goプログラムが動くための仕組み
Goプログラムを動かすためのミニOSのような存在で、以下のような機能を持つ。
| 機能              | 説明                                     |
| --------------- | -------------------------------------- |
| goroutineの管理    | goroutineを生成・スケジュールし、並列に実行できるように管理する   |
| スケジューラ          | goroutineを効率的にCPUに割り当てる（G-M-Pモデル：後述）   |
| ガベージコレクション      | 使わなくなったメモリを自動で解放する                     |
| チャネルとselect文    | goroutine間の通信を安全に管理する                  |
| タイマーやSleep      | `time.Sleep` や `time.After` などのタイミング制御 |
| panic/recover機能 | エラー処理や例外的状況への対応                        |

goroutineのスケジューリングはGoランタイムに任せているため、goroutineの実行順序は保証されていない。

# fork-joinモデル
main関数はmainゴルーチンであり、それを親として子ゴルーチンを生成する。
完了すれば親に合流する。gitみたいだね！

# 合流ポイントが無いとき
go hello()で実行したとき、親が先に終了してしまったらhello()が実行されない。
time.Sleepを思い浮かべるかもしれないが、時間は分からないのでアンチパターンである。方法の一つとしてsyncパッケージのWaitGroupがある。
```
func main() {
    var wg sync.WaitGroup
    hello := func() {
        defer wg.Done()
        fmt.Println("Hello")
    }
    wg.Add(1)
    go hello()
    wg.Wait() // ① 合流ポイント
}
```
sync.WaitGroup.Wait()を使うことで、ゴルーチンが終了するまで親をブロックしておける。

# クロージャ
クロージャとは関数内で作られる関数のこと。
クロージャは定義されたスコープの変数を参照できる。
これを以下のように変数とかgoroutineで使用できる。
```
message := "hello"
hello := func(input string) {
    fmt.Println(input)
}
hello(message)
```

```
message := "hello"
go func(input string) {
    fmt.Println(input)
}(message)
```

# goroutineでのクロージャの注意点
```
var wg sync.WaitGroup
    for _, salutation := range []string{"Hoge", "Fuga", "Piyo"} {
        wg.Add(1)
        go func() {
            defer wg.Done()
            fmt.Println(salutation)
        }()
    }
    wg.Wait()
```
このような場合、forが先に終了するため、"hoge""fuga""piyo"の出力は得られず、"piyo""piyo""piyo"が出力される。
ゴルーチンがそのメモリにアクセスし続けられるようにメモリをヒープ領域へ移してくれているので、最後の値("Piyo")を参照したメモリがヒープに移される。それをfmt.Println(salutation)が参照するため、この結果になる。
なので、クロージャ内で参照するときはクロージャの引数にsalutationを入れておくのが推奨される。