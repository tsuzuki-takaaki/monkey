https://www.oreilly.co.jp/books/9784873118222/

## Go Tips

## 字句解析器(レキサー)

ソースコードを入力として受け取り、出力としてそのソースコードを表現する**トークン列**を返す(まだ抽象構文木ではない)

字句解析器の仕事は「コードが意味をなすか・動作するか・エラーを含むか」を判定すること<font color=red>**ではない**</font>

=> **入力をトークン列に直す**ことに専念するもの

### REPL
```Read-Eval-Print-Loop```
コンソールとかインタラクティブモードとか言ったりもする(シェル的なそれ)

## 構文解析器(parser)
やっていることはJavascriptの```Json.parse```と本質的に同じ
```Json.parse```とかは、シリアライズ形式のパーサー(デシリアライズ: 文字列 -> javascriptのobject)

### parser generator
↑ その名の通り、parserを**作る**もの(parser自体ではない)

RubyKaigiでよく聞いた, ```yacc, bison```はこれに当たる

https://ja.wikipedia.org/wiki/%E3%83%91%E3%83%BC%E3%82%B5%E3%82%B8%E3%82%A7%E3%83%8D%E3%83%AC%E3%83%BC%E3%82%BF

- **言語を作る場合は自分で構文解析器を書くなんてことはほとんどなく、parser generatorを使って構文解析器を生成するのが普通**
- が、このtutorialでは構文解析器を自作する

### 構文解析器の形式
- トップダウン構文解析器: ASTのルートノードから構築を開始して下がっていく
- ボトムダウン構文解析器: ↑の逆

### ソースコードを評価するにあたって
式と文の違いは何か
- ```式```: 値を生成する(```5```) <- expression
- ```文```: 値を生成しない(```let x = 5```)

>In programming, an expression is a value, or anything that executes and ends up being a value. 

- 何が式で何が文かはその言語によって変わってくる
    - 条件分岐が式になるものもある(実際にrubyはそう)
```ruby
  a = if 10 > 9
        true
      else
        false
      end
```

### 構文解析器がやることを超ざっくり
**繰り返しトークンを読み進め、現在のトークンを調べて次にすることを決める.他の構文解析関数を呼ぶか、エラーを発生させるかのどっちか.**

## Goのドメイン知識

- Goの変数宣言
- 初期値なし
```go
    var s string
```
- 初期値あり
```go
    var s string = "hello world"
    // or
    s := "hello world"
```
- ```文字列```
    - ```ダブルクォート```: 改行を含めることはできず、バックスラッシュ\によるエスケープは解釈される。
    - ```バッククォート```: バッククォート以外の改行を含めたすべての文字を書くことができる。バックスラッシュのエスケープは解釈されない。
- ```[]Type```: 配列
    - https://qiita.com/k-penguin-sato/items/daad9986d6c42bdcde90
    - ex: ```[]int{1, 2, 3, 4}```
        - 1, 2, 3, 4を要素にもつint型の配列
- 構造体の定義方法
```go
    type [型(構造体)の名前] struct {
        [フィールド名] [型名]
        [フィールド名] [型名]
        [フィールド名] [型名]
    }
```
- golang の 引数、戻り値、レシーバをポインタにすべきか
    - goでは値のコピーが起こるからポインタを使うことによって、コピーを発生させない？
    - structとかArrayとかオーバーヘッドが大きいからポインタの方が良さそう？
    - https://www.pospome.work/entry/2017/08/12/195032
- goにおける```*```と```&```
    - https://qiita.com/tmzkysk/items/1b73eaf415fee91aaad3
    - ```&```: 変数 -> pointer
    - ```*```: pointer -> 変数
        - +αでポインタ**<font color=red>型</font>**の宣言をする時にも使う
        - ↑ アホほど紛らわしい
- Goにおけるメソッドと関数
    - https://qiita.com/yosuke_takeuchi/items/5bd061b4a766c43c9995
    - メソッドの場合はメソッド名の前に謎の引数がある(おそらくというか絶対レシーバ)
- Goのmap
    - ```var 変数名 map[キーの型]値の型```
    - ex. ```var keywards = map[string]TokenType```
- GoのカンマOK
```go
    // keywordsはmap
    if tok, ok := keywards[ident]; ok {
		return tok
	}
	return IDENT
```
↑ mapに該当のkeyがあった場合はif文がtrueとして処理される
- Goのinterface
    - https://trap.jp/post/1445/
    - ↓ 定義方法
```go
    type 名前 interface {
        メソッド メソッドでreturnされる型
    }

    // ex
    type Node interface {
        TokenLiteral() string
    }
```
↑ の例だと、「```Node```というのは、```TokenLiteralメソッド```を持つインターフェイスです」の意になる
<=> ```TokenLiteral```メソッドを持つということは```Node```インターフェイスを持っている

- <font color=red>**ダックタイピング**</font>味が強いなと思った
- Goのtype function
    - 第一級関数として、変数に入れられるため変数の型として、関数の型を定義することができる
```go
    // よくあるやつ
    var s string
    s = "hello world"
    
    // type function
    type StStFunc func(st string) string
    var func_variable StStFunc
    // ↑ このように定義すると、func_variableにはStStFunc型のデータしか入れることができなくなる
```

https://qiita.com/rock619/items/db44507d02814e490902
