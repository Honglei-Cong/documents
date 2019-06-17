
正则表达式是程序设计中经常用到的，通常被用来检索、替换那些符合某个模式的字符串。
当前大多程序设计语言，为了方便程序开发，都提供了正则表达式的标准库。
但是本体平台的智能合约目前还没有较好的正则表达式的标准库。

本次挑战为完成一个简单规则的正则表达式匹配的智能合约。需要的实现的匹配规则如下：

```
    c    matches any literal character c
    .    matches any single character
    *    matches zero or more occurrences of the previous character
    ^    matches the beginning of the input string
    $    matches the end of the input string
```

比如模式 x.y 能匹配 xay, x2y等,但它不能匹配 xy 或 xaby。
^.$ 能够与任何单个字符的字符串匹配，而 ^\*$ 能够与任意字符串匹配。

智能合约模版为：

```
def Main(regex, args):
    return match(regex, args[0])


def match(regex, target):
    // Your Implementation Here
    return True
```


示例测试用例：

1. "abc" ~ "abc"
2. "a.c" ~ "abc"
3. "a\*c" ~ "abbbbc"
4. "^\*c" ~ "abbbbc"
5. "^$" ~ ""


评测结果标准介绍：

我们准备了200个测试用例，所有字符皆为ascii字符，最长字符串长度为1024。

1. 通过所有测试用例
2. 完成所有测试的所需要的NeoVM指令数总和最少


