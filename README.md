# 💬 WasaText

**WasaText** è un'applicazione di messaggistica istantanea web-based, sviluppata come progetto universitario per il corso di "Web and Software Architecture" (WASA). L'applicazione permette agli utenti di scambiarsi messaggi testuali e immagini in tempo reale, gestire conversazioni singole e di gruppo.


## 🚀 Funzionalità Principali

L'applicazione offre un set completo di funzionalità social:

* **Autenticazione:** Login semplificato (basato su username).
* **Chat 1-to-1:** Messaggistica privata diretta tra utenti.
* **Gruppi:**
    * Creazione di gruppi.
    * Aggiunta di partecipanti.
    * Modifica nome e foto del gruppo.
    * Possibilità di abbandonare i gruppi.
* **Messaggi Multimediali:** Supporto per l'invio di immagini (GIF, PNG, JPG).
* **Interazioni Avanzate:**
    * Risposte ai messaggi (Reply).
    * Inoltro messaggi (Forward).
    * Eliminazione messaggi (Delete).
    * Reazioni/Commenti con emoji.
* **Stato Messaggi:** Visualizzazione ricevute di lettura.
* **Profilo Utente:** Modifica del nome utente e della foto profilo personalizzata.

## 🛠️ Tecnologie Utilizzate

Il progetto è costruito utilizzando un'architettura a microservizi containerizzata.

### Backend
* **Linguaggio:** [Go (Golang)](https://go.dev/)
* **Router:** `julienschmidt/httprouter`
* **Database:** SQLite (Embedded)
* **Design:** API RESTful

### Frontend
* **Framework:** [Vue.js 3](https://vuejs.org/)
* **Build Tool:** Yarn / Webpack
* **Stile:** CSS3 Custom & Responsive Design

### DevOps & Deployment
* **Containerization:** Docker
* **Orchestration:** Docker Compose
* **Web Server:** Nginx (per servire il frontend in produzione)

## 📋 Prerequisiti

Per eseguire il progetto in locale, assicurati di avere installato:

* [Docker](https://www.docker.com/products/docker-desktop)
* [Docker Compose](https://docs.docker.com/compose/install/)

## ⚡ Installazione e Avvio

Grazie a Docker, l'avvio dell'applicazione è immediato e non richiede l'installazione locale di Go o Node.js.

1.  **Clona il repository:**
    ```bash
    git clone [https://github.com/tuo-username/wasatext.git](https://github.com/tuo-username/wasatext.git)
    cd wasatext
    ```

2.  **Avvia i container (Backend + Frontend):**
    Esegui il comando di build e start:
    ```bash
    docker-compose up --build
    ```
    *Il flag `--build` assicura che eventuali modifiche al codice vengano ricompilate.*

3.  **Accedi all'applicazione:**
    Una volta terminato l'avvio, apri il browser ai seguenti indirizzi:

    * **Frontend (App):** [http://localhost:8080](http://localhost:8080)
    * **Backend (API):** [http://localhost:3000](http://localhost:3000)

4.  **Arrestare l'applicazione:**
    Per fermare i container e rimuovere le risorse temporanee:
    ```bash
    docker-compose down
    ```

## 📂 Struttura del Progetto

```text
wasatext/
├── cmd/                # Entry point del backend Go
├── service/            # Logica di business e API
│   ├── api/            # Handler REST (es. set-my-photo.go, set-token.go)
│   ├── database/       # Gestione DB SQLite e Query (es. database.go)
│   └── ...
├── webui/              # Codice sorgente Frontend (Vue.js)
│   ├── src/
│   │   ├── components/ # Componenti Vue (es. ListUser.vue, Conversation.vue)
│   │   ├── views/      # Pagine principali (es. HomeView.vue)
│   │   └── ...
├── Dockerfile.backend  # Configurazione build Backend
├── Dockerfile.frontend # Configurazione build Frontend (Multi-stage)
├── docker-compose.yml  # Orchestrazione servizi
└── README.md           # Questo file
