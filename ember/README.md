# ember

一个解析 My Clippings.txt 的脚本. 

起名的含义, kindle点燃,照亮,着火之意, 标注的内容就是火焰燃尽的所保留的有型物质--余烬(ember)

## 说明
> 当前支持中文版的 My Clippings的格式

> 标注输出格式固定
```
* 2019-07-03~2019-07-09(在临近标注内容相差15天以上会插入一条时间范围说明)
序号. 标注文本
    * 评论
```
## 启动
```
./ember -f 文件地址
```

## 支持功能
```
b 列出所有书名
f 根据书名查找
s 列出指定书籍id的标注内容
```
