#include "rwlock_sem.h"

const int32_t readerMax = ~(-1 << (sizeof(int32_t) * 8) - 1);

int32_t atomic_add_int32(volatile int32_t *count, int add)
{
    __asm__ __volatile__(
        "lock xadd %0, (%1);"
        : "=a"(add)
        : "r"(count), "a"(add)
        : "memory");
    return *count;
}

void RWMutex_init(RWMutex *locker)
{
    pthread_mutex_init(&locker->m, NULL);
    locker->readerWait = 0;
    locker->readerCnt = 0;
    sem_init(&locker->r_sem, 0, 1);
    sem_init(&locker->w_sem, 0, 1);
}

void RWMutex_rlock(RWMutex *locker)
{
    int r = atomic_add_int32(&locker->readerCnt, 1);
    if (r < 0)
    {
        //当前有个写锁, 读操作阻塞等待写锁释放
        sem_wait(&locker->r_sem);
    }
}

void RWMutex_rulock(RWMutex *locker)
{
    // 首先通过atomic的原子性使readerCount-1
    // 1.若readerCount大于0, 证明当前还有读锁, 直接结束本次操作
    // 2.若readerCount小于0, 证明已经没有读锁, 但是还有因为读锁被阻塞的写锁存在
    int32_t r = atomic_add_int32(&locker->readerCnt, -1);
    if (r < 0)
    {
        // 尝试唤醒被阻塞的写锁
        if (atomic_add_int32(&locker->readerWait, -1) == 0)
        {
            sem_post(&locker->w_sem);
        }
    }
}

void RWMutex_lock(RWMutex *locker)
{
    pthread_mutex_lock(&locker->m);
    int32_t readerCnt = atomic_add_int32(&locker->readerCnt, -readerMax) + readerMax;
    // 当r不为0说明，当前写锁之前有读锁的存在
    // 修改下readerWait，也就是当前写锁需要等待的读锁的个数
    if (readerCnt != 0 && atomic_add_int32(&locker->readerWait, readerCnt) != 0)
    {
        // 阻塞当前写锁
        sem_wait(&locker->w_sem);
    }
}

void RWMutex_ulock(RWMutex *locker)
{
    // 增加readerCount, 若超过读锁的最大限制, 触发panic
    // 和写锁定的-rwmutexMaxReaders，向对应
    int32_t r = atomic_add_int32(&locker->readerCnt, readerMax);

    // 如果r>0，说明当前写锁后面，有阻塞的读锁
    // 然后，通过信号量一一释放阻塞的读锁
    for (int32_t i = 0; i < r; i++)
    {
        sem_post(&locker->r_sem);
    }
    // 释放互斥锁
    pthread_mutex_unlock(&locker->m);
}