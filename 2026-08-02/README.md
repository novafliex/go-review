# 2026-08-02｜Go 并发模式

## 当天练习

生产者消费者、扇出、扇入、Pipeline。

**和什么有关：** channel 所有权、结果汇总。

## 小项目结构

examples/ 模式示例；exercises/ 练习测试。

## 快速复习与运行

```powershell
go test -race ./2026-08-02/exercises
```

扇出是多人处理，扇入是多路汇总。

## 常见问题

扇入要等生产者退出后再关闭输出 channel。
