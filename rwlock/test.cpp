#include <iostream>
#include <windows.h>
#include "rwlock_sem.h"
int32_t num = 0;

RWMutex rwlocker;

void *reader(void *)
{
    for (int32_t i=1;i<=100;i++)
    {
        RWMutex_rlock(&rwlocker);
        printf("thread id:%d ,read: %d\n", pthread_self(), num);
        RWMutex_rulock(&rwlocker);
    }
}
void *writer(void *)
{
    for (int32_t i=1;i<=100;i++)
    {
        RWMutex_lock(&rwlocker);
        num+=i;
        printf("thread id:%d ,write: %d\n", pthread_self(), num);
        RWMutex_ulock(&rwlocker);
    }
}

int main(int, char **)
{
    RWMutex_init(&rwlocker);

    pthread_t func1;
    pthread_create(&func1, NULL, reader, NULL);
    pthread_t func2;
    pthread_create(&func2, NULL, writer, NULL);
    // pthread_t func3;
    // pthread_create(&func3, NULL, writer, NULL);
    pthread_join(func2, NULL);
    std::cout << "Hello, world!\n";
}
