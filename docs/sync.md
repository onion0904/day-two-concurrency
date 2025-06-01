# sync パッケージとは

並行処理の同期や相互排他ロックなどの機能が備わっている。

# WaitGroup

以下のように使用

```
var wg sync.WaitGroup
wg.Add(1) // 1
    go func() {
        defer wg.Done() //2
        fmt.Println("2nd goroutine sleeping...")
        time.Sleep(2)
    }()
wg.Wait() //3
```

Add()でタスクの追加。(カウンタ++)
Done()で終了を伝える。(カウンタ--)
Wait()で WaitGroup のカウンタが 0 になるまで親をブロック

# Mutex

たとえば、トイレを思い浮かべてください。
トイレは 1 つしかない（＝共有リソース）
誰かが使っている間は、他の人は入れない（＝排他的に使う）
この「今誰かが使っているかどうかを管理する鍵が Mutex です。
そして、クリティカルセクション(同時に使用されると困る共有リソース)を扱うときに Mutex を使用する。
以下のように使用

```
func main() {
    var count int
    var lock sync.Mutex

    increment := func() {
        lock.Lock()
        defer lock.Unlock()
        count++
        fmt.Printf("increment: %d\n", count)
    }

    decrement := func() {
        lock.Lock()
        defer lock.Unlock()
        count--
        fmt.Printf("decrement %d\n", count)
    }

    // Increment
    var arithmetic sync.WaitGroup
    for i := 0; i <= 5; i++ {
        arithmetic.Add(1)
        go func() {
            defer arithmetic.Done()
            increment()
        }()
    }

    // Decrement
    for i := 0; i <= 5; i++ {
        arithmetic.Add(1)
        go func() {
            defer arithmetic.Done()
            decrement()
        }()
    }

    arithmetic.Wait()

    fmt.Println("arithmetic complete")
}
```

count 変数をゴルーチンで共有してるが、Mutex を使用することで「インクリメント時にはインクリメントだけする」といった相互排他制御が簡単に実装できます。

- Unlock を呼び出し忘れるとデッドロックするので注意
- クリティカルセクションをできるだけ短かくすること

# RWMutex

RWMutex は、メモリ管理の機能を提供してくれています。 これは、共有メモリへの同時読み込みを許容しつつ、書き込みを排他的に制御するために使用されます。例えば、キャッシュや設定データなど、頻繁に読み込まれるが稀にしか更新されないデータに対して有効です。
Mutex と同じように使用します

```
d.RLock()
defer d.RUnlock()
```

# Cond

条件変数と呼ばれ、特定の条件が満たされるのを待つために使われる。
以下のように使用

```
func main() {
    var wg sync.WaitGroup
    var dataReady bool
    cond := sync.NewCond(&sync.Mutex{}) //1

    // データが準備されるのを待つゴルーチン
    waitForData := func(i int) {
        defer wg.Done()
        cond.L.Lock()
        for !dataReady {
            fmt.Printf("ゴルーチン%d: データを待っています\n", i)
            cond.Wait() //2
        }
        fmt.Printf("ゴルーチン%d: データが準備されました\n", i)
        cond.L.Unlock()
    }

    // 5つのゴルーチンを起動し、データが準備されるのを待つ
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go waitForData(i)
    }

    // データを準備する
    go func() {
        fmt.Println("データの準備")
        time.Sleep(3 * time.Second) // データ準備に3秒かかるとする
        cond.L.Lock()
        dataReady = true
        cond.Broadcast() // 3
        cond.L.Unlock()
    }()

    wg.Wait()
}
```

cond := sync.NewCond(&sync.Mutex{})で条件変数を作成。
ゴルーチンが Wait()を呼び出すと、別のゴルーチンが同じ sync.Cond 変数で Signal()または Broadcast()を呼び出すまで実行が中断される。
Signal()はランダムに一つ、Broadcast()は全ての goroutine を起動する。

# Once
以下のようにonce.DO()の引数に関数を取り、複数のgoroutineや関数で呼ばれても一回しか実行しない
```
var once sync.Once
once.Do(func() {
		fmt.Println("once")
	})
```

# Pool
以下のように使用
```
func main() {
    myPool := &sync.Pool{
        New: func() interface{} {
            fmt.Println("Creating new instance")
            return struct{}{}
        },
    }

    myPool.Get() //1
    instance := myPool.Get() // 2
    myPool.Put(instance) //3
    myPool.Get() //4
}
```
1. Get()で初期インスタンスを取得する。変数に割り当てられてないのでガベージコレクタに捨てられる。
2. Get()で取得、変数instanceに割り当てられる
3. PutでPoolに返却
4. Get()でPoolから取得。既存のインスタンスなので、「Creating new instance」と出力されない。

- リソースのやり取りは全てany(interface{})型で行う
- sync.Poolをインスタンス化するときは、スレッド安全なNewメンバー変数を用意する