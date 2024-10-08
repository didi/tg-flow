# tg-flow简介
tg-flow是一个专注于在线高并发场景的工作流引擎，它综合运用了多种优化算法，可以在高并发场景下轻松地支持由100多个节点组成的复杂工作流的调度，目前已经广泛运用于滴滴内部多个日流量数十亿的核心在线高并发系统。

tg-flow 相关的开源仓库总共将包括3个部分:
* tg-flow:    工作流引擎核心模块，提供高并发场景下的复杂工作流的调度执行能力。
* tg-example: 为便于用户快速上手而提供的一个示例项目，如果用户本地有安装golang环境，下载后直接go run main.go即可运行。[点此了解 tg-example](https://github.com/didi/tg-example)
* tg-service: 后台管理系统，辅助用户对工作流进行配置管理及发布上线，此部分非tg-flow的强依赖模块，尚在整理中，待完善后再提交。

# 我们的目标
受限于在线高并发场景的高性能要求，传统工作流引擎很难支持在线高并发系统。 我们的目标是实现一个既能够满足在线高并发场景的性能要求，又具备传统工作流引擎各种复杂功能的工作流引擎。

# 产品特色
## 性能特色：
* tg-flow综合使用了多种高效的算法，可以在单机环境下轻松支持100多个串并行节点组成的复杂工作流的调度，除去业务逻辑之外的引擎本身的耗时仅0.1ms左右。
## 功能特色：
* 适用面广：可适用于所有在线高并发场景，典型的如搜索引擎、广告引擎、推荐引擎等。目前已经在滴滴内部多个日流量数十亿的核心在线系统中使用。
* 功能强大：
  1. 支持不限层次嵌套的串、并行混合的工作流。
  2. 支持条件分支，支持多（N>2）条件枚举分支，支持分支节点中带动态（工作流执行过程中动态赋值的）参数。
  3. 支持不限层级的父子工作流嵌套。
  4. 支持对工作流内部节点和分支进行各种超时控制，如：节点内部超时控制、节点间的超时控制、子工作流节点的超时控制，支持节点超时回调。
  5. 提供基于工作流引擎的在线问题诊断能力，可以通过在用户界面点击工作流节点显示节点中间输出结果。
* 集成灵活：
  1. 如果只想使用tg-flow的调度功能，这种情况下只需要自己按工作流DSL语法完成工作流配置文件的编写，然后放入自己在线系统的目录中，在线系统进行工作流引擎初始化时加载上述配置文件即可。
  2. 进一步地，如果你想使用tg-flow的工作流配置管理功能，可以照着用户手册搭建web服务，在web界面上以拖拽方式构建工作流，然后在系统管理页导出工作流配置文件压缩包，解压后放入自己的在线系统的目录中。
  3. 再进一步，如果需要使用tg-flow所有功能，可以在第二步的基础上，自行部署redis或zookeeper，然后在线系统中可以使用tg-core中的核心模块定期从redis或zookeeper中检测工作流新版本，并及时更新到在线系统。
  
# 用户手册
   [用户手册](https://github.com/didi/tg-example/blob/main/user_manual.md)
  
# 欢迎加入
  请联系：[张云锋](https://github.com/dayunzhangyunfeng), [周子纯](https://github.com/zhouzichun0315), [唐桂尧](https://github.com/tgy931)
