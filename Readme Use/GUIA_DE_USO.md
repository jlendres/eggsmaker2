# 📘 Guía de Uso - Eggsmaker 2

## 📋 Índice

1. [Introducción](#introducción)
2. [Instalación](#instalación)
3. [Primer Inicio](#primer-inicio)
4. [Interfaz Principal](#interfaz-principal)
5. [Flujo de Trabajo](#flujo-de-trabajo)
6. [Opciones Avanzadas](#opciones-avanzadas)
7. [Resolución de Problemas](#resolución-de-problemas)
8. [Preguntas Frecuentes](#preguntas-frecuentes)

---

## 🎯 Introducción

**Eggsmaker 2** es una aplicación gráfica para **Penguins' Eggs** que permite crear ISOs remasterizables y administrar instalaciones Live y paquetes del sistema.

### ¿Qué hace Eggsmaker 2?

- ✅ Crea ISOs de tu sistema Linux actual
- ✅ Clona tu escritorio y configuraciones
- ✅ Personaliza la configuración de la ISO
- ✅ Copia ISOs generadas a ubicaciones personalizadas
- ✅ Ejecuta la secuencia de creación ISO de Eggsmaker Plus
- ✅ Inicia instalaciones Live con Calamares o Krill
- ✅ Recopila listas de paquetes y permite instalar o desinstalar paquetes nativos, AUR, Chaotic-AUR y Flatpak
- ✅ Edita la configuración de Eggs desde un formulario gráfico
- ✅ Permite limpiar temporales al finalizar

### Requisitos Previos

- **Penguins' Eggs** instalado en el sistema
- Permisos de **sudo** (se solicitará la contraseña al inicio)
- Python 3.8 o superior
- Conexión a internet (para instalar paquetes o descargar componentes)

### Instalación Live y paquetes

La sección **Instalación Live** permite iniciar Calamares o Krill. Calamares se abre como aplicación gráfica; Krill requiere una terminal gráfica compatible y solicita allí la contraseña de `sudo` para su interfaz de texto.

En **Gestor de paquetes**, pulsa **Recopilar** para crear listas en `~/$(hostname)-pak/`, selecciona un origen y los paquetes, y usa **Instalar selección** o **Desinstalar selección**. Los orígenes disponibles dependen de la distribución y de las herramientas instaladas. En Arch se admiten paquetes nativos, AUR mediante `yay`, Chaotic-AUR y Flatpak; en Debian/Ubuntu se gestionan paquetes manuales de APT y Flatpak.

El botón **Configurar Eggs** abre el editor gráfico de `eggs.yaml`; no ejecuta el asistente interactivo `eggs config` en modo texto.

---

## 📦 Instalación

### Método 1: Instalador Automático (Recomendado)

#### Instalación Local (Usuario)
```bash
cd /ruta/a/eggsmaker2
./install.sh
```

Esto instalará Eggsmaker 2 en `~/.local/share/eggsmaker2`

#### Instalación Global (Sistema)
```bash
cd /ruta/a/eggsmaker2
sudo ./install.sh --system
```

Esto instalará Eggsmaker 2 en `/opt/eggsmaker2`

### Método 2: Instalación con pip

```bash
cd /ruta/a/eggsmaker2
pip install -e .
```

### Método 3: Ejecución Directa (Sin Instalación)

```bash
cd /ruta/a/eggsmaker2
python3 main.py
```

---

## 🚀 Primer Inicio

### 1. Ejecutar la Aplicación

Después de la instalación:
```bash
eggsmaker2
```

O desde el menú de aplicaciones: buscar **"Eggsmaker 2"**

### 2. Autenticación Inicial

Al iniciar, se mostrará un diálogo solicitando la **contraseña de sudo**:

```
┌─────────────────────────────┐
│ Autenticación requerida     │
│ (sudo)                      │
│                             │
│ Contraseña de sudo: [****] │
│                             │
│ [Aceptar]  [Cancelar]      │
└─────────────────────────────┘
```

> 💡 **Importante**: Esta contraseña se requiere para ejecutar comandos privilegiados de Penguins' Eggs.

### 3. Navegador Web

La aplicación se abrirá automáticamente en tu navegador predeterminado en:
```
http://localhost:8080
```

Para acceder desde otro dispositivo, inicia la aplicación con `EGGSMAKER_HOST=0.0.0.0 eggsmaker2` y usa la dirección indicada en el pie de página. Limita el acceso a redes confiables.

---

## 🖥️ Interfaz Principal

### Barra Superior

```
┌──────────────────────────────────────────────────────────────┐
│ 🥚 Eggsmaker - Versión 26.10.21      Eggsmaker              │
│                                                              │
│                              Penguins' Eggs: X.X.X           │
└──────────────────────────────────────────────────────────────┘
```

### Terminal de Salida

```
┌──────────────────────────────────────────────────────────────┐
│ Terminal de salida (log de comandos y operaciones)         │
│                                                              │
│ > Iniciando...                                              │
│ > ...                                                        │
└──────────────────────────────────────────────────────────────┘
```

### Acciones principales

```
Crear ISO  |  Instalación Live  |  Configurar Eggs  |  Gestor de paquetes
```

### Barra de Estado

```
┌──────────────────────────────────────────────────────────────┐
│ Estado: Esperando                                            │
│ ████████████████████░░░░░░░░░ 65%                          │
│                                                              │
│ 65% | Copias: 2 | ISO: 3.5 GB | Copia: 00:02:30 | ...      │
└──────────────────────────────────────────────────────────────┘
```

---

## 🔄 Flujo de Trabajo

1. En la pregunta de configuración del escritorio, elige si ejecutar `sudo eggs tools skel`.
2. Si eliges editar la configuración de Eggs, se abre el formulario gráfico de `eggs.yaml`; no se lanza el asistente de texto `eggs config`.
3. Elige si quieres crear una ISO normal o un clon del sistema (`--clone`).
4. Confirma la generación. Eggsmaker ejecuta `sudo eggs kill && eggs tools clean && eggs remaster`, añadiendo `--clone` cuando corresponde.
5. Selecciona el destino en el cuadro de diálogo. El selector de carpetas de la aplicación se reutiliza aquí y `rsync --info=progress2` muestra la copia en el terminal y la barra de progreso.
6. Puedes repetir la copia en otros destinos o elegir omitirla.
7. Al final puedes ejecutar `eggs kill` para quitar los temporales de construcción.

El orden es secuencial y cada decisión se presenta gráficamente antes de ejecutar la siguiente acción.

## ⚡ Opciones Avanzadas

### 🔧 Editar Configuración ISO

Disponible mediante **Configurar Eggs**, que abre un formulario gráfico con los campos frecuentes y el YAML completo.

#### Campos Editables:

| Campo | Descripción | Ejemplo |
|-------|-------------|---------|
| **Contraseña de root** | Contraseña del superusuario en la ISO | `toor123` |
| **Nombre base del snapshot** | Nombre de la distribución | `mi-debian` |
| **Prefijo del snapshot** | Prefijo para versiones | `custom-` |
| **Contraseña de usuario** | Contraseña del usuario live | `live123` |

> 📝 **Archivo editado**: `/etc/penguins-eggs.d/eggs.yaml`

### 📊 Indicadores en Tiempo Real

#### Cronómetros

| Cronómetro | Color | Descripción |
|------------|-------|-------------|
| **Copia** | Cyan | Tiempo de copia actual |
| **Generación** | Rojo | Tiempo de generación ISO |
| **Total** | Verde Lima | Suma de todos los tiempos |

#### Barra de Progreso

- **Azul claro**: Operación activa
- **Ondulante**: Generación de ISO (duración indeterminada)
- **Roja**: Copia de ISO (muestra porcentaje exacto)

#### Indicadores Adicionales

- **Tamaño ISO**: Se actualiza después de generar
- **Copias realizadas**: Contador de copias exitosas
- **Porcentaje**: Durante la copia (0-100%)

---

## 🔍 Resolución de Problemas

### ❌ Error: "Se requiere contraseña de sudo"

**Causa**: No se proporcionó la contraseña al inicio o expiró la sesión.

**Solución**: Recarga la página web y proporciona la contraseña nuevamente.

---

### ❌ Error: "No se encontró ISO para copiar"

**Causa**: La ISO no se generó correctamente o no está en la ubicación esperada.

**Soluciones**:
1. Verifica que la generación terminó correctamente
2. Revisa el terminal de log para errores durante la generación
3. Verifica manualmente si existe el archivo:
   ```bash
   ls -lh /home/eggs/*.iso
   ls -lh /home/eggs/.mnt/*.iso
   ```

---

### ⚠️ La aplicación no abre en el navegador

**Soluciones**:
1. Abre manualmente: `http://localhost:8080`
2. Verifica que el puerto esté disponible:
   ```bash
   netstat -tuln | grep 8080
   ```
3. Usa un puerto diferente:
   ```bash
   PORT=8081 eggsmaker2
   ```

---

### ⚠️ "Permiso denegado" al copiar ISO

**Causa**: La ISO está en un directorio protegido.

**Solución**: La aplicación automáticamente intentará usar `sudo` para copiar.

---

### 🐛 El terminal no muestra salida

**Solución**: Recarga la página web (F5).

---

## ❓ Preguntas Frecuentes

### ¿Cuánto espacio necesito para crear una ISO?

Depende del modo:
- **Sin incluir datos**: ~2-5 GB
- **Con datos (/home completo)**: Tamaño de tu `/home` + 2-5 GB
- **Máxima compresión**: ~70-80% del tamaño normal

> 💡 **Consejo**: Usa "Limpiar sesión" después de copiar la ISO para liberar espacio.

---

### ¿Puedo usar Eggsmaker 2 desde otro dispositivo en la red?

Por seguridad, la aplicación escucha solo en `127.0.0.1` de forma predeterminada. Para habilitar la red local, inicia con:
```
EGGSMAKER_HOST=0.0.0.0 eggsmaker2
```

El pie de página mostrará la dirección de acceso. Usa esta opción solo en una red confiable: la interfaz puede ejecutar acciones administrativas después de autenticar `sudo`.

---

### ¿Qué sucede si cierro el navegador durante una operación?

La operación **continúa ejecutándose** en segundo plano. Puedes:
1. Volver a abrir `http://localhost:8080`
2. Ver el progreso en el terminal

---

### ¿Puedo cancelar una operación en progreso?

Actualmente no hay botón de cancelación. Opciones:
1. Esperar a que termine
2. Cerrar la aplicación completamente (Ctrl+C en la terminal)
3. Usar "Limpiar sesión" después para limpiar archivos parciales

---

### ¿Dónde se guardan las ISOs generadas?

Ubicaciones típicas (verificadas en orden):
1. `/home/eggs/` - Ubicación principal
2. `/home/eggs/.mnt/` - Ubicación alternativa

El tamaño se muestra automáticamente en la barra de estado.

---

### ¿Puedo crear múltiples ISOs sin limpiar sesión?

Sí, pero:
- ⚠️ La nueva ISO sobrescribirá la anterior en `/home/eggs/`
- 💡 **Consejo**: Copia cada ISO antes de generar la siguiente

---

### ¿Cómo actualizo Eggsmaker 2?

```bash
cd /ruta/a/eggsmaker2
git pull
./install.sh  # Reinstalar
```

O si usaste pip:
```bash
pip install -e . --upgrade
```

---

## 📚 Recursos Adicionales

### Enlaces Útiles

- **Penguins' Eggs**: https://penguins-eggs.net/
- **Repositorio GitHub**: https://github.com/pieroproietti/penguins-eggs
- **Documentación Eggs**: https://penguins-eggs.net/docs/

### Soporte

Para problemas específicos de Eggsmaker 2, consulta:
- El terminal de log en la aplicación
- Los archivos de log del sistema
- La comunidad de Penguins' Eggs

---

## 📝 Resumen

1. Pulsa **Crear ISO**.
2. Responde si quieres copiar la configuración del escritorio, editar `eggs.yaml` y crear un clon.
3. Confirma la generación y sigue su progreso en el terminal y la barra.
4. Selecciona una carpeta con el selector integrado, repite la copia o sáltala.
5. Elige si ejecutar `eggs kill` para quitar temporales.

## 🎉 ¡Disfruta creando tus ISOs!

Esta guía cubre todas las funcionalidades principales de Eggsmaker 2. Para usos avanzados, consulta la documentación de Penguins' Eggs.

**Versión de la guía**: 1.0
**Compatible con**: Eggsmaker 2 v26.10.21
