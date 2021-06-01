#ifndef RWLOCK_H
#define RWLOCK_H

#include <pthread.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C"{
#endif

typedef struct RWMutex_s
{
    pthread_mutex_t m;
    pthread_cond_t rcond;
    pthread_cond_t wcond;
    int32_t readerCnt;
    int32_t writerCnt;
}RWMutex;

void RWMutex_inti(RWMutex* locker);
void RWMutex_rlock(RWMutex* locker);
void RWMutex_rulock(RWMutex* locker);
void RWMutex_lock(RWMutex* locker);
void RWMutex_ulock(RWMutex* locker);
#ifdef __cplusplus
}
#endif
#endif