Exercise 1
### 3 Sharing a variable
c:

> Oppretter funksjoner som skal utføres av tråden. Inni main oppretter vi trådene ved bruk av thread_create, sørger for at dette går vha. if. Bruker join som gjør at tråden får kjørt ferdig før programmet går videre. Da vi kjørte programmmet før bruk av join, fikk vi lavere resultater for i, sannsynligvis fordi tråden ikke fikk kjørt seg ferdig. Etter bruk av join: høyere (mer positiv) / lavere (mer negativ) verdi for i. (NULL inni join sin if er peker til exitstatus - forklrer hvorfor exit)

go:

> Oppretter funksjoner her også. Disse utføres ved hjelp av Goroutine. Det er basically en thread i Go, altså det som kaller på funksjonen, og får funksjonen til å kjøre. GOMAXPROCS setter maks antall prosesser (her tråder) som kan kjøres samtidig. Sette til 1--> kan bare kjøre en tråd. Goroutine virker som at kjører raskt, for når vi setter sleep til lave verdier får vi like store tall som ved store. Sette sleep til 0 gir 0, altså får ikke goroutine kjørt funksjonen. 

### 4 Sharing a variable, but properly
 Semaphore: everyone can unlock
 Mutex: only the owner has the key
 Vi ønsker å unngå race conditions og data race, slik vi fikk i oppgave 3.

c:
 - For å unngå race conditions og data race er MUTEX best fordi vi slipper at to tråder forsøker å skrive til samme minneområde samtidig.
 - Bruker lock, og ikke trylock, fordi lock venter på sin tur til å låse når det blir ledig, mens trylock returenerer istedenfor å vente.
 - Fungerer!!  Også med minus 1.
 
go:
> Vi brukte select med cases. Lagde fire channels for sending av inc, dec, read og quit, eget signal for hver handling. Lagde en server som kan endre på i. Funksjonene signaliserer riktig handling med å sende et signal over riktig kanal. Har lagt til en quit-lytte-funksjonalitet som sender over egen quit-kanal. Denne gjør at vi venter på at decrement og increment gjør seg ferdig. Select går alltid tilbake der den var hvis den blir avbrutt midt i en handling. Dette gjør at vi alltid kommer til null, fordi alle oppgaver blir utført. 

### 5 Bounded buffer
c:

> Initialiserte to semaforer for bufferen som skal lages. Disse har initialverdi lik size og 0, altså en teller opp, og en teller ned. Disse brukes for å sørge for fornufit push- og pop-funksjonalitet for stacken. Semaforer er counters som fungerer slik:
- post increments by one
- wait decrements semaphore by one
- cannot decrement if semaphore is zero -> blocks
La til mutex rundt push og pop, for å sørge for at disse operasjonene skjer fullstendig/ uten avbrytelse.

go:
> Bruker Channel buffer. Denne introduserer en slgs "kanalkapasitet", altså hvor mange elementer som kan være "på kanalen" til enhver tid. Det gjør at vi kan fortsette å skrive til kanalen uten å tømme den, så lenge kapasiteten ikke er fyllt opp. Bare implementere make med et ekstra tall. 
