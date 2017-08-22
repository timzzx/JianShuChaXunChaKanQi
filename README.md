# 简书mini查询阅读器
用go walk写的一个简书的关键字查询来查看的阅读器
  **开始** 
![输入图片说明](https://git.oschina.net/uploads/images/2017/0821/193913_9a81ce53_462123.jpeg "aaaaaa.jpg")
 **使用查询** 
![输入图片说明](https://git.oschina.net/uploads/images/2017/0821/194308_7e1ffe9d_462123.jpeg "bbb.jpg")


**运行需要** 
- golang版本需要1.8以上
- go get github.com/lxn/walk
- go get github.com/akavel/rsrc
  
**执行** 
- go build -ldflags="-H windowsgui"

**或者** 
- go build
 
**说明：** 
- 学习golang1个月，一直在看看语法，写写hello world。前两天突然想做个桌面应用，正好在浏览简书，就决定做个简书的查看器。
- 花了一天看了walk的demo就直接开写。golang我只会皮毛，walk的资料也太少了只能看源码。费了点力气完成的这个简陋的软件。
 
**问题：** 
-  1.查询只能10秒查一次，这是简书的限制，翻页不受限制。（这个问题修改不难，但是我不打算改）
-  2.没有加入登录功能（本来就是个demo你还想干啥）
-  3.代码中用了两种网页抓取的方式，net/http, goquery(为了学习)
-  4.json的解析写的很low。（毕竟刚学）
-  5.快速选择查看的标签会导致退出，原因已经找到。walk的webview的问题，网页还未加载好，重新刷新url就会导致webview关闭，webview一旦关闭就导致程序退出。（这个问题暂时没有好的解决办法）