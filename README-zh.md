# Bamboo

## 什么是 Bamboo？

**Bamboo** 是一个原型开发和评估框架，用于研究区块链专用的下一代 BFT（拜占庭容错）协议，即链式 BFT或 cBFT。
利用 Bamboo，开发人员可以在约 300 行代码中开发出全新的 cBFT 协议原型，并利用丰富的基准测试工具进行评估。

Bamboo 的设计基于这样一种观点，即 cBFT 协议的核心可以抽象为 4 条规则：**提议**、**投票**、**状态更新**和**提交**。
因此，Bamboo 将这 4 条规则抽象为一个*安全*模块，并提供可在 cBFT 协议中共享的其余组件的实现，而安全模块则由开发人员指定。

*警告*：**Bamboo**仍在大力开发中，还有更多的功能和协议有待加入。

Bamboo 的详细信息可参见本[技术报告](https://arxiv.org/abs/2103.00777)。该论文将发表于 [ICDCS 2021](https://icdcs2021.us/)。

## 什么是 cBFT？

在高层次上，cBFT 协议共享一个统一的 "提议-投票" 范式，即在全局分类账中为来自客户端的交易分配一个唯一的顺序。
区块链是通过哈希值加密连接在一起的区块序列。
区块链中的每个区块都包含其父区块的哈希值以及一批交易和其他元数据。

与经典的 BFT 协议类似，cBFT 协议由领导节点驱动，以逐个视图的方式运行。
每个参与者在收到信息后都会根据协议的四条特定规则采取行动：**提议**、**投票**、**状态更新**和**提交**。
每个视图都有一个随机选择的指定领导者，该领导者根据**提议**规则提议一个区块，并填充网络。
收到区块后，副本根据**投票**规则采取行动，并根据**状态更新**规则更新其本地状态。
对于每个视图，副本都应通过为区块创建*法定证书*（或 QC）来证明所提议区块的有效性。
具有有效 QC 的区块即被视为经过认证。
区块链的基本结构如下图所示。

![blockchain](https://github.com/gitferry/bamboo/blob/master/doc/propose-vote.jpeg?raw=true)

分叉发生的原因是区块冲突，即两个区块不能相互扩展。
出现区块冲突的原因可能是网络延迟或提议者故意忽略区块链尾部。
只要区块满足基于其本地状态的**提交**规则，副本就会最终确定该区块。

一旦一个区块最终确定，整个区块链的前缀也会最终确定。规则规定，所有最终确定的区块必须保留在一条链中。

最终确定的区块可从内存中移除，转入持久化存储空间，以便进行垃圾回收。

## 本仓库包含哪些内容？

协议：

- [x] [HotStuff 和 two-chain HotStuff](https://dl.acm.org/doi/10.1145/3293611.3331591)
- [x] [Streamlet](https://dl.acm.org/doi/10.1145/3419614.3423256)
- [x] [Fast-HotStuff](https://arxiv.org/abs/2010.11454)
- [ ] [LBFT](https://arxiv.org/abs/2012.01636)
- [ ] [SFT](https://arxiv.org/abs/2101.03715)

功能：

- [x] 基准测试
- [x] 故障注入

## 如何编译

1. 安装 [Go](https://golang.org/dl/).

2. 下载 Bamboo 源码。

3. 编译 `server` 和 `client`。

```sh
cd bamboo/bin
go build ../server
go build ../client
```

## 如何运行

用户可以在模拟模式（单进程）或部署模式中运行基于 Bamboo 的 cBFT 协议。

### 模拟模式

在模拟模式下，副本在不同的 Go routines 中运行，信息通过 Go channels 传递。

1. `cd bamboo/bin`
2. 修改 `ips.txt` 中每个节点的 IP 地址集。IP 数等于节点数。此处，本地 IP 为 `127.0.0.1`。每个节点将被分配一个从 `8070` 起递增的端口。
3. 修改 `config.json` 中的配置参数。
4. 修改 `simulation.sh` 以指定要运行的协议名称。
5. 运行 `server` ，然后使用脚本运行 `client`。

```sh
bash simulation.sh
bash runClient.sh
```

1. 依次停止客户端和服务端，关闭模拟。

```sh
bash closeClient.sh
bash stop.sh
```

日志生成在本地目录，名称为 `client/server.xxx.log`，其中 `xxx` 是进程的 pid。

### 部署模式

Bamboo 可以部署在真实网络中。

1. `cd bamboo/bin/deploy`
2. 编译 `server` 和 `client`。
3. 分别在 `pub_ips.txt` 和 `ips.txt` 中指定服务端节点的外部 IP 和内部 IP。
4. 在 `clients.txt` 中指定作为客户机运行的机器的 IP 地址。
5. 在 `run.sh` 中指定协议类型。
6. 修改 `config.json` 中的配置参数。
7. 修改 `deploy.sh` 和 `setup_cli.sh` 以指定登录服务端和客户端的用户名和密码。
8. 将二进制文件和配置文件上传到远程机器上。

```sh
bash deploy.sh
bash setup_cli.sh
```

1. 在远程机器上上传/更新配置文件。

   ```sh
   bash update_conf.sh
   ```

2. 启动服务端节点。

   ```sh
   bash start.sh
   ```

3. 通过 ssh 登录并启动客户端（假设只有一台）。

   ```sh
   bash ./runClient.sh
   ```

   可以在 `runClient.sh` 中指定并发客户端的数量。

4. 停止客户端和服务端。

   ```sh
   bash ./closeClient.sh
   bash ./pkill.sh
   ```

### 监视器

在每次运行期间，可以通过浏览器查看节点的统计数据（吞吐量、延迟、视图编号等）。

http://127.0.0.1:8070/query

其中 `127.0.0.1:8070` 可以替换为实际的节点地址。
