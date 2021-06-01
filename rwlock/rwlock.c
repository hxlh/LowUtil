#include "rwlock.h"

int32_t atomic_add_int32(volatile int32_t *count, int add)
{
    __asm__ __volatile__(
        "lock xadd %0, (%1);"
        : "=a"(add)
        : "r"(count), "a"(add)
        : "memory");
    return *count;
}

void RWMutex_inti(RWMutex *locker)
{
    pthread_mutex_init(&locker->m, NULL);
    pthread_cond_init(&locker->rcond, NULL);
    pthread_cond_init(&locker->wcond, NULL);
    locker->writerCnt = 0;
    locker->readerCnt = 0;
    
}

void RWMutex_rlock(RWMutex *locker) 
{
    pthread_mutex_lock(&locker->m);
    while (locker->writerCnt != 0)
    {
        pthread_cond_wait(&locker->rcond, &locker->m); //等待条件变量的成立
    }
    locker->readerCnt++;
    pthread_mutex_unlock(&locker->m);
}

void RWMutex_rulock(RWMutex *locker)
{
    int32_t r=atomic_add_int32(&locker->readerCnt, -1);
    if (r == 0)
    {
        //唤醒写线程
        pthread_cond_broadcast(&locker->wcond);
    }
}

void RWMutex_lock(RWMutex *locker)
{
    pthread_mutex_lock(&locker->m); //加锁
    while (locker->writerCnt > 0 || locker->readerCnt>0)
    {
        pthread_cond_wait(&locker->wcond, &locker->m); //等待条件变量的成立
    }
    locker->writerCnt++;
}

void RWMutex_ulock(RWMutex *locker)
{
    locker->writerCnt--;
    if(locker->writerCnt==0)
    {
        //唤醒读进程
        pthread_cond_broadcast(&locker->rcond);
    }
    pthread_mutex_unlock(&locker->m);
}