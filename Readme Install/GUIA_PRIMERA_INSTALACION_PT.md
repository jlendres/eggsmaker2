# Guia do Usuário: Primeira Instalação e Portabilidade

Este guia detalha os arquivos necessários para portar este projeto para outros computadores e como realizar a instalação.

## 1. Arquivos Necessários para Portar

Para instalar este projeto em outros computadores, você deve copiar a pasta completa do projeto. Certifique-se de que ela inclua os seguintes arquivos e pastas essenciais:

### Código Fonte e Recursos
*   **`main.py`**: Arquivo principal da aplicação.
*   **`backend.py`**: Lógica do sistema.
*   **`version.py`**: Versão da aplicação.
*   **`assets/`**: Pasta com imagens e recursos (logo, etc.).

### Instalação
*   **`install.sh`**: Script de instalação automatizada.
*   **`requirements.txt`**: Dependências do Python.
*   **`setup.py`** e **`pyproject.toml`**: Configuração do pacote.
*   **`eggsmaker2.desktop`**: Atalho para o menu de aplicativos.

---

## 2. Instruções de Instalação

Uma vez copiada a pasta do projeto para o novo computador, você tem duas opções de instalação, dependendo das suas necessidades.

### Opção A: Instalação Local (Recomendada)
Ideal para um único usuário. Não afeta o sistema operacional globalmente.

1.  Abra um terminal dentro da pasta do projeto.
2.  Execute o seguinte comando:
    ```bash
    ./install.sh
    ```

**Localização dos arquivos:**
*   O programa será instalado em: `~/.local/share/eggsmaker2/`
*   O executável estará em: `~/.local/bin/eggsmaker2`

### Opção B: Instalação no Sistema
Torna o programa disponível para todos os usuários do computador. Requer permissões de administrador.

1.  Abra um terminal dentro da pasta do projeto.
2.  Execute o seguinte comando com `sudo`:
    ```bash
    sudo ./install.sh --system
    ```

**Localização dos arquivos:**
*   O programa será instalado em: `/opt/eggsmaker2/`
*   O executável estará em: `/usr/local/bin/eggsmaker2`

---

## 3. Execução

Uma vez instalado, você pode procurar por **"Eggsmaker 2"** no menu de aplicativos do seu sistema ou executá-lo a partir do terminal com o comando:

```bash
eggsmaker2
```

---

## 4. Atualização

Se você precisar atualizar o aplicativo para uma nova versão:

1.  Baixe ou copie a nova versão do projeto.
2.  Repita a etapa de instalação (Opção A ou B) que você usou originalmente.
    *   O instalador substituirá automaticamente os arquivos antigos pelos novos.
    *   Você não perderá suas configurações pessoais se elas estiverem fora da pasta do programa.

---

## 5. Desinstalação

Para remover o aplicativo completamente:

1.  Abra um terminal na pasta do projeto (se você ainda a tiver) ou baixe o script `uninstall.sh`.
2.  Execute o desinstalador de acordo com o seu tipo de instalação:

**Se você instalou localmente (Opção A):**
```bash
./uninstall.sh
```

**Se você instalou no sistema (Opção B):**
```bash
sudo ./uninstall.sh --system
```


