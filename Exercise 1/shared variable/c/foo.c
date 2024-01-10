// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t myLock;


// Note the return type: void*
void* incrementingThreadFunction(){
    // TODO: increment i 1_000_000 times
    
    for ( int j = 0; j < 1000000; j++)
    {   
        pthread_mutex_lock(&myLock);
        i++;
        pthread_mutex_unlock(&myLock);
    }
    
    return NULL;
}

void* decrementingThreadFunction(){
    // TODO: decrement i 1_000_000 times
    for (int j = 0; j < 1000001; j++)
    {
        pthread_mutex_lock(&myLock);
        i--;
        pthread_mutex_unlock(&myLock);
    }
    return NULL;
}


int main(){

    // TODO: 
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?

    pthread_mutex_init(&myLock, NULL);

    pthread_t thread_1;
    pthread_t thread_2;


    if (pthread_create(&thread_1, NULL, incrementingThreadFunction, "thread 1") != 0) {
    perror("pthread_create() error");
    }

    if (pthread_create(&thread_2, NULL, decrementingThreadFunction, "thread 2") !=0 ) {
    perror("pthread_create() error");
    }
    //implementere MUTEX for å få til concurrency
    
    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`   
    //join gjør at tråden får kjørt ferdig før programmet går videre
    if (pthread_join(thread_1, NULL) != 0) {perror("pthread_create() error");} //NULL er peker til exitstatus - forklrer hvorfor exit
    if (pthread_join(thread_2, NULL) != 0) {perror("pthread_create() error");}
    
    printf("The magic number is: %d\n", i);
    return 0;
}

