# Guida Utente: Prima Installazione e Portabilità

Questa guida descrive i file necessari per portare questo progetto su altri computer e come eseguire l'installazione.

## 1. File Necessari per il Porting

Per installare questo progetto su altri computer, è necessario copiare l'intera cartella del progetto. Assicurati che includa i seguenti file e cartelle essenziali:

### Codice Sorgente e Risorse
*   **`main.py`**: File principale dell'applicazione.
*   **`backend.py`**: Logica del sistema.
*   **`version.py`**: Versione dell'applicazione.
*   **`assets/`**: Cartella con immagini e risorse (logo, ecc.).

### Installazione
*   **`install.sh`**: Script di installazione automatizzata.
*   **`requirements.txt`**: Dipendenze Python.
*   **`setup.py`** e **`pyproject.toml`**: Configurazione del pacchetto.
*   **`eggsmaker2.desktop`**: Collegamento desktop per il menu.

---

## 2. Istruzioni per l'Installazione

Una volta copiata la cartella del progetto sul nuovo computer, hai due opzioni di installazione a seconda delle tue esigenze.

### Opzione A: Installazione Locale (Consigliata)
Ideale per un singolo utente. Non influisce sul sistema operativo globale.

1.  Apri un terminale all'interno della cartella del progetto.
2.  Esegui il seguente comando:
    ```bash
    ./install.sh
    ```

**Posizione dei file:**
*   Il programma verrà installato in: `~/.local/share/eggsmaker2/`
*   L'eseguibile si troverà in: `~/.local/bin/eggsmaker2`

### Opzione B: Installazione di Sistema
Rende il programma disponibile per tutti gli utenti del computer. Richiede permessi di amministratore.

1.  Apri un terminale all'interno della cartella del progetto.
2.  Esegui il seguente comando con `sudo`:
    ```bash
    sudo ./install.sh --system
    ```

**Posizione dei file:**
*   Il programma verrà installato in: `/opt/eggsmaker2/`
*   L'eseguibile si troverà in: `/usr/local/bin/eggsmaker2`

---

## 3. Esecuzione

Una volta installato, puoi cercare **"Eggsmaker 2"** nel menu delle applicazioni del tuo sistema o eseguirlo dal terminale con il comando:

```bash
eggsmaker2
```

---

## 4. Aggiornamento

Se devi aggiornare l'applicazione a una nuova versione:

1.  Scarica o copia la nuova versione del progetto.
2.  Ripeti il passaggio di installazione (Opzione A o B) utilizzato in origine.
    *   Il programma di installazione sovrascriverà automaticamente i vecchi file con quelli nuovi.
    *   Non perderai le configurazioni personali se si trovano all'esterno della cartella del programma.

---

## 5. Disinstallazione

Per rimuovere completamente l'applicazione:

1.  Apri un terminale nella cartella del progetto (se ce l'hai ancora) o scarica lo script `uninstall.sh`.
2.  Esegui il programma di disinstallazione in base al tipo di installazione:

**Se hai installato localmente (Opzione A):**
```bash
./uninstall.sh
```

**Se hai installato nel sistema (Opzione B):**
```bash
sudo ./uninstall.sh --system
```


