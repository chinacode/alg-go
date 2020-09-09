##linux state 

* R running runnable 状态
* S interrupted sleep 可中断休眠状态
* D disk sleep (uninterrupted sleep) 不可中断休眠状态，等待某些硬件程序
* I idle 空闲状态
* Z zombie 僵尸状态，进程已经运行结束，父进程未还未回收 
* T stopped || traced 表示进程初始暂停或跟踪状态（gdb调试，断点）
* X dead 死亡状态一般看不见

```go

```
