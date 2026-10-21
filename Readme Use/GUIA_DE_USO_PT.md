# 📘 Guia do Usuário - Eggsmaker 2

## 📋 Índice

1. [Introdução](#introdução)
2. [Instalação](#instalação)
3. [Primeiro Início](#primeiro-início)
4. [Interface Principal](#interface-principal)
5. [Fluxo de Trabalho](#fluxo-de-trabalho)
6. [Opções Avançadas](#opções-avançadas)
7. [Resolução de Problemas](#resolução-de-problemas)
8. [FAQ](#faq)

---

## 🎯 Introdução

**Eggsmaker 2** é um aplicativo gráfico para **Penguins' Eggs** que permite criar ISOs remasterizáveis e gerenciar instalações Live e pacotes do sistema.

### O que o Eggsmaker 2 faz?

- ✅ Cria ISOs do seu sistema Linux atual
- ✅ Clona seu desktop e configurações
- ✅ Personaliza a configuração da ISO
- ✅ Copia ISOs geradas para locais personalizados
- ✅ Executa a sequência de criação ISO do Eggsmaker Plus
- ✅ Inicia instalações Live com Calamares ou Krill
- ✅ Coleta listas de pacotes e permite instalar ou remover pacotes nativos, AUR, Chaotic-AUR e Flatpak
- ✅ Edita a configuração do Eggs em um formulário gráfico
- ✅ Permite remover arquivos temporários ao final

### Pré-requisitos

- **Penguins' Eggs** instalado no sistema
- Permissões de **sudo** (será solicitada a senha no início)
- Python 3.8 ou superior
- Conexão com a internet (para instalar pacotes ou baixar componentes)

### Instalação Live e pacotes

A seção **Instalação Live** inicia o Calamares ou o Krill. O Calamares abre como aplicativo gráfico; o Krill requer um terminal gráfico compatível e solicita nele a senha `sudo` para sua interface de texto.

Em **Gerenciador de pacotes**, use **Coletar** para criar listas em `~/$(hostname)-pak/`, escolha uma origem e os pacotes e use **Instalar seleção** ou **Desinstalar seleção**. As origens disponíveis dependem da distribuição e das ferramentas instaladas. **Configurar Eggs** abre o editor gráfico de `eggs.yaml` em vez do assistente interativo `eggs config` em modo texto.

---

## 📦 Instalação

### Método 1: Instalador Automático (Recomendado)

#### Instalação Local (Usuário)
```bash
cd /caminho/para/eggsmaker2
./install.sh
```

Isso instalará o Eggsmaker 2 em `~/.local/share/eggsmaker2`

#### Instalação Global (Sistema)
```bash
cd /caminho/para/eggsmaker2
sudo ./install.sh --system
```

Isso instalará o Eggsmaker 2 em `/opt/eggsmaker2`

### Método 2: Instalação com pip

```bash
cd /caminho/para/eggsmaker2
pip install -e .
```

### Método 3: Execução Direta (Sem Instalação)

```bash
cd /caminho/para/eggsmaker2
python3 main.py
```

---

## 🚀 Primeiro Início

### 1. Executar o Aplicativo

Após a instalação:
```bash
eggsmaker2
```

Ou no menu de aplicativos: procure por **"Eggsmaker 2"**

### 2. Autenticação Inicial

Ao iniciar, aparecerá uma caixa de diálogo solicitando a **senha de sudo**:

```
┌─────────────────────────────┐
│ Autenticação necessária     │
│ (sudo)                      │
│                             │
│ Senha sudo: [****]          │
│                             │
│ [OK]       [Cancelar]       │
└─────────────────────────────┘
```

> 💡 **Importante**: Esta senha é necessária para executar comandos privilegiados do Penguins' Eggs.

### 3. Navegador Web

O aplicativo abrirá automaticamente no seu navegador padrão em:
```
http://localhost:8080
```

Para acessar de outro dispositivo, inicie o aplicativo com `EGGSMAKER_HOST=0.0.0.0 eggsmaker2` e use o endereço mostrado no rodapé. Restrinja o acesso a redes confiáveis.

---

## 🖥️ Interface Principal

### Barra Superior

```
┌──────────────────────────────────────────────────────────────┐
│ 🥚 Eggsmaker - Versão 26.10.21       Eggsmaker              │
│                                                              │
│                              Penguins' Eggs: X.X.X           │
└──────────────────────────────────────────────────────────────┘
```

### Terminal de Saída

```
┌──────────────────────────────────────────────────────────────┐
│ Terminal de saída (log de comandos e operações)            │
│                                                              │
│ > Iniciando...                                              │
│ > ...                                                        │
└──────────────────────────────────────────────────────────────┘
```

### Ações principais

```
Criar ISO  |  Instalação Live  |  Configurar Eggs  |  Gerenciador de pacotes
```

### Barra de Status

```
┌──────────────────────────────────────────────────────────────┐
│ Status: Aguardando                                           │
│ ████████████████████░░░░░░░░░ 65%                          │
│                                                              │
│ 65% | Cópias: 2 | ISO: 3.5 GB | Cópia: 00:02:30 | ...      │
└──────────────────────────────────────────────────────────────┘
```

---

## 🔄 Fluxo de trabalho

1. Na pergunta sobre a configuração do desktop, escolha se deseja executar `sudo eggs tools skel`.
2. Se escolher editar a configuração do Eggs, será aberto o formulário gráfico de `eggs.yaml`; o assistente de texto `eggs config` não será iniciado.
3. Escolha entre criar uma ISO normal ou um clone do sistema (`--clone`).
4. Confirme a geração. O Eggsmaker executa `sudo eggs kill && eggs tools clean && eggs remaster`, acrescentando `--clone` quando selecionado.
5. Escolha o destino na janela. O seletor de pastas do aplicativo é reutilizado, e `rsync --info=progress2` mostra o progresso no terminal e na barra.
6. Você pode copiar para outros destinos ou ignorar a cópia.
7. No final, pode executar `eggs kill` para remover os arquivos temporários de compilação.

A sequência é linear; cada decisão aparece graficamente antes da próxima ação.

## ⚡ Opções Avançadas

### 🔧 Editar Configuração ISO

Disponível em **Configurar Eggs**, que abre um formulário gráfico com os campos principais e o arquivo YAML completo.

#### Campos Editáveis:

| Campo | Descrição | Exemplo |
|-------|-----------|---------|
| **Senha de root** | Senha do superusuário na ISO | `toor123` |
| **Nome base do snapshot** | Nome da distribuição | `minha-distro` |
| **Prefixo do snapshot** | Prefixo para versões | `custom-` |
| **Senha de usuário** | Senha do usuário live | `live123` |

> 📝 **Arquivo editado**: `/etc/penguins-eggs.d/eggs.yaml`

### 📊 Indicadores em Tempo Real

#### Cronômetros

| Cronômetro | Cor | Descrição |
|------------|-----|-----------|
| **Cópia** | Ciano | Tempo de cópia atual |
| **Geração** | Vermelho | Tempo de geração ISO |
| **Total** | Verde Lima | Soma de todos os tempos |

#### Barra de Progresso

- **Azul claro**: Operação ativa
- **Ondulada**: Geração de ISO (duração indeterminada)
- **Vermelha**: Cópia de ISO (mostra porcentagem exata)

#### Indicadores Adicionais

- **Tamanho ISO**: Atualiza após a geração
- **Cópias realizadas**: Contador de cópias bem-sucedidas
- **Porcentagem**: Durante a cópia (0-100%)

---

## 🔍 Resolução de Problemas

### ❌ Erro: "Senha de sudo necessária"

**Causa**: A senha não foi fornecida no início ou a sessão expirou.

**Solução**: Recarregue a página web e forneça a senha novamente.

---

### ❌ Erro: "ISO para copiar não encontrada"

**Causa**: A ISO não foi gerada corretamente ou não está no local esperado.

**Soluções**:
1. Verifique se a geração foi concluída corretamente
2. Verifique o terminal de log para erros durante a geração
3. Verifique manualmente se o arquivo existe:
   ```bash
   ls -lh /home/eggs/*.iso
   ls -lh /home/eggs/.mnt/*.iso
   ```

---

### ❌ Erro durante atualização Eggs/Calamares

**Causa**: Problemas de rede ou permissões.

**Soluções**:
1. Verifique sua conexão com a internet
2. Tente manualmente:
   ```bash
   cd ~
   git clone https://github.com/pieroproietti/fresh-eggs
   cd fresh-eggs
   sudo ./fresh-eggs.sh
   ```

---

### ⚠️ O aplicativo não abre no navegador

**Soluções**:
1. Abra manualmente: `http://localhost:8080`
2. Verifique se a porta está disponível:
   ```bash
   netstat -tuln | grep 8080
   ```
3. Use uma porta diferente:
   ```bash
   PORT=8081 eggsmaker2
   ```

---

### ⚠️ "Permissão negada" ao copiar ISO

**Causa**: A ISO está em um diretório protegido.

**Solução**: O aplicativo tentará automaticamente usar `sudo` para copiar.

---

### 🐛 O terminal não mostra saída

**Solução**: Recarregue a página web (F5).

---

## ❓ FAQ

### Quanto espaço eu preciso para criar uma ISO?

Depende do modo:
- **Sem dados**: ~2-5 GB
- **Com dados (/home completo)**: Tamanho do seu `/home` + 2-5 GB
- **Máxima compressão**: ~70-80% do tamanho normal

> 💡 **Dica**: Use "Limpar sessão" depois de copiar a ISO para liberar espaço.

---

### Posso usar o Eggsmaker 2 de outro dispositivo na rede?

Por padrão, o aplicativo só fica disponível localmente. Para habilitar o acesso pela rede, use `EGGSMAKER_HOST=0.0.0.0 eggsmaker2` e conecte-se ao endereço mostrado no rodapé. Use apenas em uma rede confiável, pois a interface pode executar ações administrativas.

---

### O que acontece se eu fechar o navegador durante uma operação?

A operação **continua em execução** em segundo plano. Você pode:
1. Reabrir `http://localhost:8080`
2. Ver o progresso no terminal

---

### Posso cancelar uma operação em andamento?

Atualmente não há botão de cancelar. Opções:
1. Esperar terminar
2. Fechar o aplicativo completamente (Ctrl+C no terminal)
3. Usar "Limpar sessão" depois para limpar arquivos parciais

---

### Onde as ISOs geradas são salvas?

Locais típicos (verificados em ordem):
1. `/home/eggs/` - Local principal
2. `/home/eggs/.mnt/` - Local alternativo

O tamanho é mostrado automaticamente na barra de status.

---

### Posso criar várias ISOs sem limpar a sessão?

Sim, mas:
- ⚠️ A nova ISO substituirá a anterior em `/home/eggs/`
- 💡 **Dica**: Copie cada ISO antes de gerar a próxima

---

### Como atualizo o Eggsmaker 2?

```bash
cd /caminho/para/eggsmaker2
git pull
./install.sh  # Reinstalar
```

Ou se você usou pip:
```bash
pip install -e . --upgrade
```

---

## 📚 Recursos Adicionais

### Links Úteis

- **Penguins' Eggs**: https://penguins-eggs.net/
- **Repositório GitHub**: https://github.com/pieroproietti/penguins-eggs
- **Documentação Eggs**: https://penguins-eggs.net/docs/

### Suporte

Para problemas específicos do Eggsmaker 2, verifique:
- O terminal de log no aplicativo
- Os arquivos de log do sistema
- A comunidade do Penguins' Eggs

---

## 📝 Resumo

1. Clique em **Criar ISO**.
2. Escolha se deseja copiar a configuração do desktop, editar `eggs.yaml` e criar um clone do sistema.
3. Confirme a geração e acompanhe o progresso no terminal e na barra.
4. Selecione uma pasta no seletor integrado, copie novamente ou ignore a cópia.
5. Escolha se deseja executar `eggs kill` para remover arquivos temporários.

## 🎉 Divirta-se criando suas ISOs!

Este guia cobre todas as principais funcionalidades do Eggsmaker 2. Para uso avançado, consulte a documentação do Penguins' Eggs.

**Versão do guia**: 1.0
**Compatível com**: Eggsmaker 2 v26.10.21
