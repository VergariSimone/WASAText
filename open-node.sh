#!/usr/bin/env sh

# Apre un ambiente di sviluppo Node.js 20 "usa e getta" (--rm).
# - Monta la cartella corrente in /src per editing live.
# - Usa 'network host' per esporre le porte direttamente su localhost.
# - Esegue come root (-u 0:0) per evitare problemi di permessi in scrittura.
# NOTA: Richiede l'esistenza della sottocartella 'webui'.

docker run -it --rm -v "$(pwd):/src:rw" -u "0:0" --network host --workdir /src/webui node:20 /bin/bash