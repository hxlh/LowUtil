#ifndef RWLOCK_SEM_H
#define RWLOCK_SEM_H

#include <pthread.h>
#include <semaphore.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C"{
#endif

//参考go读写锁实现 https://www.cnblogs.com/ricklz/p/13697322.html

typedef struct RWMutex_s
{
    pthread_mutex_t m;
    sem_t r_sem;
    sem_t w_sem;
    int32_t readerCnt;
    int32_t readerWait;
}RWMutex;

void RWMutex_init(RWMutex* locker);
void RWMutex_rlock(RWMutex* locker);
void RWMutex_rulock(RWMutex* locker);
void RWMutex_lock(RWMutex* locker);
void RWMutex_ulock(RWMutex* locker);

#ifdef __cplusplus
}
#endif
#endif