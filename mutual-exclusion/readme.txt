passo 1: gerar os executáveis
    $ go build Process.go
    $ go build SharedResource.go

passo 2: executar os processos e o SharedResource
    $ ./Process 1 :10002 :10003 :10004
    $ ./Process 2 :10002 :10003 :10004
    $ ./Process 3 :10002 :10003 :10004