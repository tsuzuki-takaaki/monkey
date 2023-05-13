https://www.oreilly.co.jp/books/9784873118222/

## Go Tips

## 字句解析器(レキサー)

ソースコードを入力として受け取り、出力としてそのソースコードを表現する**トークン列**を返す(まだ抽象構文木ではない)

字句解析器の仕事は「コードが意味をなすか・動作するか・エラーを含むか」を判定すること<font color=red>**ではない**</font>

=> **入力をトークン列に直す**ことに専念するもの

## Goのドメイン知識

- ```:=```: 初期化かつ代入
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

https://qiita.com/rock619/items/db44507d02814e490902
