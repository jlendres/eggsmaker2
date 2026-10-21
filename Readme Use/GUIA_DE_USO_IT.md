# 📘 Guida Utente - Eggsmaker 2

## 📋 Indice

1. [Introduzione](#introduzione)
2. [Installazione](#installazione)
3. [Primo Avvio](#primo-avvio)
4. [Interfaccia Principale](#interfaccia-principale)
5. [Flusso di Lavoro](#flusso-di-lavoro)
6. [Opzioni Avanzate](#opzioni-avanzate)
7. [Risoluzione dei Problemi](#risoluzione-dei-problemi)
8. [FAQ](#faq)

---

## 🎯 Introduzione

**Eggsmaker 2** è un'applicazione grafica per **Penguins' Eggs** che consente di creare ISO rimasterizzabili e gestire installazioni Live e pacchetti di sistema.

### Cosa fa Eggsmaker 2?

- ✅ Crea ISO del tuo sistema Linux attuale
- ✅ Clona il tuo desktop e le configurazioni
- ✅ Personalizza la configurazione della ISO
- ✅ Copia le ISO generate in posizioni personalizzate
- ✅ Esegue la sequenza di creazione ISO di Eggsmaker Plus
- ✅ Avvia installazioni Live con Calamares o Krill
- ✅ Raccoglie elenchi di pacchetti e consente di installare o rimuovere pacchetti nativi, AUR, Chaotic-AUR e Flatpak
- ✅ Modifica la configurazione di Eggs tramite un modulo grafico
- ✅ Consente di rimuovere i file temporanei al termine

### Prerequisiti

- **Penguins' Eggs** installato nel sistema
- Permessi di **sudo** (verrà richiesta la password all'avvio)
- Python 3.8 o superiore
- Connessione internet (per installare pacchetti o scaricare componenti)

### Installazione Live e pacchetti

La sezione **Installazione Live** avvia Calamares o Krill. Calamares si apre come applicazione grafica; Krill richiede un terminale grafico compatibile e chiede lì la password `sudo` per la sua interfaccia testuale.

In **Gestione pacchetti**, usa **Raccogli** per creare gli elenchi in `~/$(hostname)-pak/`, seleziona una sorgente e i pacchetti, quindi scegli **Installa selezione** o **Disinstalla selezione**. Le sorgenti dipendono dalla distribuzione e dagli strumenti installati. **Configura Eggs** apre l'editor grafico di `eggs.yaml` invece della procedura interattiva `eggs config` in modalità testo.

---

## 📦 Installazione

### Metodo 1: Installazione Automatica (Consigliato)

#### Installazione Locale (Utente)
```bash
cd /percorso/a/eggsmaker2
./install.sh
```

Questo installerà Eggsmaker 2 in `~/.local/share/eggsmaker2`

#### Installazione Globale (Sistema)
```bash
cd /percorso/a/eggsmaker2
sudo ./install.sh --system
```

Questo installerà Eggsmaker 2 in `/opt/eggsmaker2`

### Metodo 2: Installazione con pip

```bash
cd /percorso/a/eggsmaker2
pip install -e .
```

### Metodo 3: Esecuzione Diretta (Senza Installazione)

```bash
cd /percorso/a/eggsmaker2
python3 main.py
```

---

## 🚀 Primo Avvio

### 1. Eseguire l'Applicazione

Dopo l'installazione:
```bash
eggsmaker2
```

O dal menu delle applicazioni: cerca **"Eggsmaker 2"**

### 2. Autenticazione Iniziale

All'avvio, apparirà una finestra di dialogo che richiede la **password di sudo**:

```
┌─────────────────────────────┐
│ Autenticazione richiesta    │
│ (sudo)                      │
│                             │
│ Password sudo: [****]       │
│                             │
│ [OK]       [Annulla]        │
└─────────────────────────────┘
```

> 💡 **Importante**: Questa password è necessaria per eseguire comandi privilegiati di Penguins' Eggs.

### 3. Browser Web

L'applicazione si aprirà automaticamente nel tuo browser predefinito all'indirizzo:
```
http://localhost:8080
```

Per accedere da un altro dispositivo, avvia l'applicazione con `EGGSMAKER_HOST=0.0.0.0 eggsmaker2` e usa l'indirizzo mostrato nel piè di pagina. Limita l'accesso alle reti fidate.

---

## 🖥️ Interfaccia Principale

### Barra Superiore

```
┌──────────────────────────────────────────────────────────────┐
│ 🥚 Eggsmaker - Versione 26.10.21     Eggsmaker              │
│                                                              │
│                              Penguins' Eggs: X.X.X           │
└──────────────────────────────────────────────────────────────┘
```

### Terminale di Output

```
┌──────────────────────────────────────────────────────────────┐
│ Terminale di output (log comandi e operazioni)             │
│                                                              │
│ > Avvio in corso...                                         │
│ > ...                                                        │
└──────────────────────────────────────────────────────────────┘
```

### Azioni principali

```
Creazione ISO  |  Installazione Live  |  Configura Eggs  |  Gestione pacchetti
```

### Barra di Stato

```
┌──────────────────────────────────────────────────────────────┐
│ Stato: In attesa                                             │
│ ████████████████████░░░░░░░░░ 65%                          │
│                                                              │
│ 65% | Copie: 2 | ISO: 3.5 GB | Copia: 00:02:30 | ...       │
└──────────────────────────────────────────────────────────────┘
```

---

## 🔄 Flusso di lavoro

1. Alla domanda sulla configurazione del desktop, scegli se eseguire `sudo eggs tools skel`.
2. Se scegli di modificare la configurazione di Eggs, si apre il modulo grafico `eggs.yaml`; non viene avviata la procedura testuale `eggs config`.
3. Scegli se creare una ISO normale o un clone del sistema (`--clone`).
4. Conferma la generazione. Eggsmaker esegue `sudo eggs kill && eggs tools clean && eggs remaster`, aggiungendo `--clone` se selezionato.
5. Scegli la destinazione nella finestra. Il selettore di cartelle dellapplicazione viene riutilizzato e `rsync --info=progress2` mostra il progresso nel terminale e nella barra.
6. Puoi copiare in altre destinazioni oppure saltare la copia.
7. Alla fine puoi eseguire `eggs kill` per rimuovere i file temporanei di compilazione.

La sequenza è lineare e ogni decisione viene presentata graficamente prima dellazione successiva.

## ⚡ Opzioni Avanzate

### 🔧 Modifica Configurazione ISO

Disponibile tramite **Configura Eggs**, che apre un modulo grafico con i campi principali e l'intero file YAML.

#### Campi Modificabili:

| Campo | Descrizione | Esempio |
|-------|-------------|---------|
| **Password di root** | Password del superutente nella ISO | `toor123` |
| **Nome base snapshot** | Nome della distribuzione | `mia-debian` |
| **Prefisso snapshot** | Prefisso per versioni | `custom-` |
| **Password utente** | Password utente live | `live123` |

> 📝 **File modificato**: `/etc/penguins-eggs.d/eggs.yaml`

### 📊 Indicatori in Tempo Reale

#### Cronometri

| Cronometro | Colore | Descrizione |
|------------|--------|-------------|
| **Copia** | Ciano | Tempo di copia attuale |
| **Generazione** | Rosso | Tempo di generazione ISO |
| **Totale** | Verde Lime | Somma di tutti i tempi |

#### Barra di Avanzamento

- **Azzurro**: Operazione attiva
- **Ondulata**: Generazione ISO (durata indeterminata)
- **Rossa**: Copia ISO (mostra percentuale esatta)

#### Indicatori Aggiuntivi

- **Dimensione ISO**: Si aggiorna dopo la generazione
- **Copie effettuate**: Contatore di copie riuscite
- **Percentuale**: Durante la copia (0-100%)

---

## 🔍 Risoluzione dei Problemi

### ❌ Errore: "Password di sudo richiesta"

**Causa**: La password non è stata fornita all'avvio o la sessione è scaduta.

**Soluzione**: Ricarica la pagina web e fornisci nuovamente la password.

---

### ❌ Errore: "ISO da copiare non trovata"

**Causa**: La ISO non è stata generata correttamente o non si trova nella posizione prevista.

**Soluzioni**:
1. Verifica che la generazione sia terminata correttamente
2. Controlla il terminale di log per errori durante la generazione
3. Verifica manualmente se il file esiste:
   ```bash
   ls -lh /home/eggs/*.iso
   ls -lh /home/eggs/.mnt/*.iso
   ```

---

### ❌ Errore durante aggiornamento Eggs/Calamares

**Causa**: Problemi di rete o permessi.

**Soluzioni**:
1. Controlla la tua connessione internet
2. Prova manualmente:
   ```bash
   cd ~
   git clone https://github.com/pieroproietti/fresh-eggs
   cd fresh-eggs
   sudo ./fresh-eggs.sh
   ```

---

### ⚠️ L'applicazione non si apre nel browser

**Soluzioni**:
1. Apri manualmente: `http://localhost:8080`
2. Verifica che la porta sia disponibile:
   ```bash
   netstat -tuln | grep 8080
   ```
3. Usa una porta diversa:
   ```bash
   PORT=8081 eggsmaker2
   ```

---

### ⚠️ "Permesso negato" durante la copia della ISO

**Causa**: La ISO si trova in una directory protetta.

**Soluzione**: L'applicazione tenterà automaticamente di usare `sudo` per copiare.

---

### 🐛 Il terminale non mostra output

**Soluzione**: Ricarica la pagina web (F5).

---

## ❓ FAQ

### Quanto spazio mi serve per creare una ISO?

Dipende dalla modalità:
- **Senza dati**: ~2-5 GB
- **Con dati (/home completo)**: Dimensione della tua `/home` + 2-5 GB
- **Massima compressione**: ~70-80% della dimensione normale

> 💡 **Consiglio**: Usa "Pulisci sessione" dopo aver copiato la ISO per liberare spazio.

---

### Posso usare Eggsmaker 2 da un altro dispositivo nella rete?

Per impostazione predefinita, l'applicazione è disponibile solo in locale. Per abilitarla nella rete, usa `EGGSMAKER_HOST=0.0.0.0 eggsmaker2` e connettiti all'indirizzo mostrato nel piè di pagina. Usala solo su una rete fidata, poiché può eseguire operazioni amministrative.

---

### Cosa succede se chiudo il browser durante un'operazione?

L'operazione **continua l'esecuzione** in background. Puoi:
1. Riaprire `http://localhost:8080`
2. Vedere l'avanzamento nel terminale

---

### Posso annullare un'operazione in corso?

Attualmente non c'è un pulsante di annullamento. Opzioni:
1. Attendere che finisca
2. Chiudere completamente l'applicazione (Ctrl+C nel terminale)
3. Usare "Pulisci sessione" dopo per pulire i file parziali

---

### Dove vengono salvate le ISO generate?

Posizioni tipiche (controllate in ordine):
1. `/home/eggs/` - Posizione principale
2. `/home/eggs/.mnt/` - Posizione alternativa

La dimensione viene mostrata automaticamente nella barra di stato.

---

### Posso creare più ISO senza pulire la sessione?

Sì, ma:
- ⚠️ La nuova ISO sovrascriverà la precedente in `/home/eggs/`
- 💡 **Consiglio**: Copia ogni ISO prima di generare la successiva

---

### Come aggiorno Eggsmaker 2?

```bash
cd /percorso/a/eggsmaker2
git pull
./install.sh  # Reinstallare
```

O se hai usato pip:
```bash
pip install -e . --upgrade
```

---

## 📚 Risorse Aggiuntive

### Link Utili

- **Penguins' Eggs**: https://penguins-eggs.net/
- **Repository GitHub**: https://github.com/pieroproietti/penguins-eggs
- **Documentazione Eggs**: https://penguins-eggs.net/docs/

### Supporto

Per problemi specifici di Eggsmaker 2, consulta:
- Il terminale di log nell'applicazione
- I file di log del sistema
- La comunità di Penguins' Eggs

---

## 📝 Riepilogo

1. Premi **Crea ISO**.
2. Scegli se copiare la configurazione del desktop, modificare `eggs.yaml` e creare un clone.
3. Conferma la generazione e controlla il progresso nel terminale e nella barra.
4. Seleziona una cartella con il selettore integrato, ripeti la copia o saltala.
5. Scegli se eseguire `eggs kill` per rimuovere i file temporanei.

## 🎉 Divertiti a creare le tue ISO!

Questa guida copre tutte le funzionalità principali di Eggsmaker 2. Per usi avanzati, consulta la documentazione di Penguins' Eggs.

**Versione della guida**: 1.0
**Compatibile con**: Eggsmaker 2 v26.10.21
