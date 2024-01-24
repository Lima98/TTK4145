Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency means we can have several processes going at the same time in contrast to parallelism which only has one thing happening at any given moment.

What is the difference between a *race condition* and a *data race*? 
> Race condition refers to when the timing or ordering of events in a program affects its ability to perform the tasks correctly. Data race however is then two different threads access the same part of memory with non-read operations and are not synchronized.
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> It assures that things happens in a controlled manner, that resources are shared correclty, by setting resources as in use and making other processe wait for their turn to access them.


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> This is to be able to have several things running at once independently of one another. This solves the problem of having to program everything "linearly".

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> Fibers are threads of threads. They are smaller sub-routines which use co-operative multitasking in contrast to preemptive multitasking used by threads. The fiber have to cooperate in order to work they do this by yielding at different times to allow the other to run.

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> This makes our lives easier, it makes the code more readable, when you know how it works. It might take a bit more work to get this going, but when you are using them, debugging and understanding the sturcture of the program gets easier.

What do you think is best - *shared variables* or *message passing*?
> I would think shared variable are best since they are working with the same problem directly, no copies of one another, I suspec this is a more effective way of working, however I also think that it can lead to more problems since they data is shared compared to separate in message passing.


