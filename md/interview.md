### Golang相关

#### 1. GMP原理与调度

1. 单进程带来的两个问题

> 1. 单一的执行流程，计算机只能一个任务一个任务处理
> 2. 进程阻塞所带来的CPU时间浪费 

2. 多进程，多线程并发，有了调度器需求

- 带来新的问题 

> 1. 高内存占用
> 2. 调度的高消耗CPU 

3. 协程跟线程是有区别的，线程由CPU调度是抢占式的，协程由用户态调度是协作式的，一协程让出CPU后才执行下一个协程

4. Go协程 

P -- Processor。处理器

M -- 线程。

G -- goroutine 协程

P有本地队列，跟M是p：m = 1:1，本地队列存放的G数量不能超过256个，G优先加入到P的本地队列

- 调度器的设计策略
  1. work stealing机制(当本线程无G时，尝试从其他线程绑定的P偷取G，而不是销毁线程)
  2. hand off机制(当本线程因为G进行系统调度阻塞时，线程会释放绑定的P，把P转移给其他空闲的线程执行)

#### 2. Go build -ldflags参数的用处

- -s的作用是去掉符号信息。-w的作用是去掉调试信息，-X注入importantPath.name=value，importantPath是路径，name是变量名称，value是值

### MySql相关

#### 1. B树和B+树的区别

- B树：每个节点都存储key和data，所有节点组成这棵树，并且叶子结点指针为null
- B+树：只有叶子结点存储data，叶子结点包含了这棵树的所有键值，叶子结点不存储指针，后来在B+树上增加了顺序访问指针，也就是每个叶子结点增加一个指向相邻叶子节点的指针，这样一棵树成了数据库系统实现索引的首选数据结构

#### 2. mysql优化

- 设计：存储引擎，字段类型，范式与逆范式
- 功能：索引，缓存，分区分表
- 架构：主从复制，读写分离，负载均衡
- 合理SQL：测试，经验

#### 3.Mysql如何定位慢查询

- 根据慢日志定位慢查询sql

  > Show variables like '%query%';  查询慢日志相关信息

  主要看这三个属性：

  	 	-  long_query_time: 10.000000 查询超过10秒定义为慢查询
  	 	-  slow_query_log: OFF 是否打开慢查询日志
  	 	-  Slow_query_log_file: /usr/local/mysql/data/slow.sql  慢查询文件所在位置

  使用以下命令设置这些属性值：

  > set global slow_query_log = ON; #打开慢查询日志
  >
  > set global long_query_time.= 1; #超过1秒的语句被定义为慢查询， 注意设置了之后需要重新连接才有效

  使用以下命令查询慢查询日志的数量：

  > show status like '%slow_queries%';

  要查看慢日志，直接打开slow.sql文件查看即可

  > sudo vim /usr/local/mysql/data/slow.sql

- 使用explain等工具分析sql

### Redis相关

#### 1.Redis为什么这么快

- Redis完全基于内存，绝大部分请求是纯粹的内存操作，非常迅速，数据存在内存中，类似于HashMap，HashMap的优势就是查找和操作的时间复杂度都是O(1)
- 数据结构简单，对数据的操作也简单
- 采用单线程，避免了不必要的上下文切换和竞争条件，不存在多线程导致的CPU切换，不用去考虑各种锁的问题，不存在加锁释放锁操作，没有死锁问题导致的性能消耗
- 使用多路复用IO模型，非阻塞IO 

#### 2.Redis与Memcached的区别

#### 3. Redis五种数据类型的应用场景

- String : 缓存功能，计数器，共享session，限流
- List： 消息队列，最新列表，排行榜
- Hash：用户信息，购物车管理等
- Set：求交集，并集，差集，随机数（抽奖功能），社交需求（共同好友）
- Zset：实时排行榜，延时队列，限流（滑动窗口）



### K8S相关和Docker相关

#### 1.删除所有不用的images

> docker image prune -a -f

#### 2.Docker资源隔离类型

- 主要有文件系统，网络，进程pid，权限（用户和用户组）

### Kafka相关

### HTTP，TCP/IP相关

#### 1.SYN攻击查看方式

> netstat -nap | grep SYN_REVC

#### 2.实现负载均衡方案

- 随机算法，轮询算法，加权轮询算法，hash算法，一致性hash算法

#### 3. TCP/IP协议

- 应用层 -> 传输层 -> 网络层 -> 链路层

### 操作系统相关

#### 1.进程挂起恢复

- 命令最后加上& :表示将进程放到后台执行
- jobs： 查看后台进程
- fg: 恢复后台一个进程到前台执行， + 优先级高。制定某一个进程恢复到前台来。用jobs查看后台的进程编号

\- d l c b p s 文件的类型 -rwxr--r-- 


### Go-mirco全链路追踪
- https://github.com/asim/go-micro/tree/master/plugins/wrapper/trace



crucial industry certain check month arm comic elbow parent noise champion turkey
