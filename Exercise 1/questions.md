Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> *Your answer here*
>Parallellism: flere ting skjer akkurat samtidig, men dette krever flere steder for operasjon
>Concurrency: flere ting skjer "samtidig", men dette betyr at vi i praksis bare bytter veldig raskt mellom hvilke oppgaver som utføres.
>Ulike minneområder vs samme minneområde for ulike operasjoner. Evt helt separate maskiner vs en maskin. Flere CPUer vs en CPU som bytter på hvem som har tilgang veldig raskt 

What is the difference between a *race condition* and a *data race*? 
> *Your answer here*
>Race condition: Hvis du prøver å bruke data som ikke er ferdig modifisert kan du få race conditions. Hvis timing eller rekkefølge ikke hånteres riktig.
>Data race: Hvis du ikke har god kontroll på hvem som har lov til å skrive data til et minneområde til enhver tid, kan du få problemer av typen data race. 
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> *Your answer here* 
>Passer på at oppgaver blir utført i riktig rekkefølge. Scheduler selects among runnable threads. Så det vi har gjort er å gi scheduleren begrensninger på hvilke tråder som er runnable. 

### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> *Your answer here*
>Gjør at vi får concurrency - at vi kan kjøre flere oppgaver "samtidig", se beskrivelse i spm 1 av samtidig, swapping.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> *Your answer here*
>Deling av tråder i mindre tråder
>Fibers are threads of threads. They are smaller sub-routines which use co-operative multitasking in contrast to preemptive multitasking used by threads. The fiber have to cooperate in order to work they do this by yielding at different times to allow the other to run.

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> *Your answer here*
> This makes our lives easier, it makes the code more readable, when you know how it works. It might take a bit more work to get this going, but when you are using them, debugging and understanding the sturcture of the program gets easier.

What do you think is best - *shared variables* or *message passing*?
> *Your answer here*
>Shared variables: Bruker 
>Message passing: Channels
>Message passsing is best for avoiding concurrency issues, but these can also be handlede in good ways when using shared variables. So both may be used, but which is best?
