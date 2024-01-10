Exercise 1
>3 Sharing a variable
c:
Oppretter funksjoner som skal utføres av tråden. Inni main oppretter vi trådene ved bruk av thread_create, sørger for at dette går vha. if. Bruker join som gjør at tråden får kjørt ferdig før programmet går videre. Da vi kjørte programmmet før bruk av join, fikk vi lavere resultater for i, sannsynligvis fordi tråden ikke fikk kjørt seg ferdig. Etter bruk av join: høyere (mer positiv) / lavere (mer negativ) verdi for i. (NULL inni join sin if er peker til exitstatus - forklrer hvorfor exit)

go:
