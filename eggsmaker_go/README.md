# Eggsmaker Go

Aplicación gráfica para Linux basada en Penguins Eggs. Esta versión conserva el flujo principal de Eggsmaker Plus: copia opcional de configuración de escritorio, edición interactiva de `eggs config`, creación de ISO normal o clon, copia de ISO, limpieza de temporales, instaladores Calamares/Krill y gestión de listas de paquetes.

## Requisitos

- Go 1.23 o posterior.
- Penguins Eggs instalado y disponible como `eggs`.
- Debian/Ubuntu/BigLinux o Arch/CachyOS/Manjaro/EndeavourOS.
- Dependencias gráficas de Fyne/OpenGL (el instalador intenta comprobarlas).
- `pkexec` recomendado para autorización gráfica; si no está disponible, ejecutar `sudo -v` antes de iniciar Eggsmaker para autorizar la sesión.
- Una terminal gráfica para el asistente interactivo `eggs config` y Krill.

## Compilar e instalar

Para compilar manualmente (Go 1.23+ y dependencias gráficas de Fyne instaladas):

```bash
./compilar.sh
```

El script aplica `gofmt`, ejecuta `go test ./...` y genera el binario `eggsmaker`.

Para instalar en Arch Linux o Debian y derivados:

```bash
chmod +x install.sh uninstall.sh compilar.sh
./install.sh
```

El instalador comprueba Go y las bibliotecas de desarrollo de OpenGL/X11, compila como usuario normal y solicita `sudo` únicamente para instalar el ejecutable, el icono y el acceso del menú (`eggsmaker.desktop`).

La aplicación debe iniciarse como usuario normal, no con `sudo`:

```bash
eggsmaker
```

Desinstalación (el script solicita permisos de administrador si hace falta):

```bash
./uninstall.sh
```

## Generar AppImage (portable)

En Linux x86_64, con Go 1.23+, GCC y las dependencias de desarrollo de Fyne instaladas, ejecutá:

```bash
./compilar_appimage.sh
```

El script ejecuta las pruebas, compila Eggsmaker y utiliza `linuxdeploy` + `appimagetool` para recopilar las bibliotecas compartidas necesarias. Las herramientas de empaquetado se descargan automáticamente en `build/appimage/tools/` y se reutilizan en las siguientes compilaciones. El resultado queda en `dist/Eggsmaker_Go-<version>-x86_64.AppImage`. Se requiere conexión a Internet la primera vez para descargar las herramientas.

El AppImage no incluye Penguins Eggs ni los gestores de paquetes del sistema: `eggs`, `pacman`, `apt` y las herramientas externas necesarias deben estar instaladas en el equipo anfitrión. Tampoco incluye el kernel, controladores gráficos ni componentes centrales de glibc; las bibliotecas gráficas del sistema siguen siendo necesarias.

Para ejecutar el archivo generado:

```bash
chmod +x dist/Eggsmaker_Go-*-x86_64.AppImage
./dist/Eggsmaker_Go-*-x86_64.AppImage
```

La primera versión del empaquetador está orientada a x86_64. La compatibilidad final debe comprobarse en las distribuciones de destino, especialmente en sistemas con bibliotecas gráficas antiguas.

## Idiomas

La interfaz está disponible en español (predeterminado), inglés, portugués e italiano. El idioma se cambia desde el selector de la ventana principal y la preferencia queda guardada para los siguientes inicios.

## Apariencia y colores

El botón **Apariencia** (arriba a la derecha, junto al selector de idioma) permite configurar los colores de toda la aplicación. Los cambios se aplican al instante y se guardan para los siguientes inicios.

- **Tema**: *Predeterminado (actual)* conserva la apariencia original de Eggsmaker; *Seguir al sistema* adapta todos los colores al tema claro u oscuro del sistema (y cambia en vivo si el sistema lo hace); *Oscuro* y *Claro* fijan el tema sin importar el sistema.
- **Color de acento**: naranja (predeterminado), azul, verde, violeta, rojo, cian o un color personalizado elegido con el selector. Afecta a títulos, botón activo, barras de progreso, marco y texto de la terminal integrada.
- **Restablecer valores predeterminados** vuelve a la apariencia original (tema predeterminado + naranja).

## Funciones

- `eggs config`: se abre en una terminal interactiva; los cambios se guardan desde el propio asistente de Penguins Eggs. En el flujo de generación se pregunta si se desea abrirlo.
- ISO: `eggs kill && eggs tools clean && eggs remaster` o `eggs remaster --clone`.
- Progreso ISO: Penguins Eggs no proporciona un porcentaje global estable para `eggs remaster`; se muestra una barra indeterminada violeta y la salida real del proceso. La barra de copia es naranja intenso y porcentual.
- Copia: porcentaje calculado a partir de bytes escritos y tamaño de archivo. El selector de destino se abre ampliado. Se copia primero a un temporal y luego se renombra para evitar dejar un archivo final incompleto.
- Terminal integrada: visor de salida estándar y errores en tiempo real, texto cyan y marco gris claro. `eggs config` y Krill usan terminal externa porque necesitan una TTY interactiva.
- Instaladores: `eggs sysinstall calamares` y `eggs sysinstall krill`.
- Gestor: recopila listas de paquetes nativos, AUR, Chaotic-AUR y Flatpak cuando están disponibles; permite seleccionar paquetes para restaurar o desinstalar; las ventanas de selección se amplían para mostrar aproximadamente diez elementos o más sin desplazarse. La restauración debe hacerse en una distribución de la misma familia y con los repositorios/herramientas correspondientes configurados.

## Seguridad

- No se formatean ni borran particiones.
- Las operaciones de instalación/desinstalación muestran confirmación antes de ejecutarse.
- Los nombres de paquetes se pasan como argumentos separados y se filtran antes de ejecutar gestores de paquetes.
- Eggsmaker Plus no forma parte de este proyecto y no debe modificarse.

## Notas de compatibilidad

- La interfaz se ejecuta como usuario normal. Las tareas administrativas se elevan individualmente mediante `pkexec`; si no está instalado, se usa `sudo -n` y conviene ejecutar `sudo -v` antes de abrir la aplicación.
- La terminal externa permanece abierta hasta que se presiona Enter, de modo que se pueden revisar los mensajes de `eggs config` y guardar los cambios dentro del asistente.
- La barra de progreso de copia es porcentual y se actualiza con el tamaño real transferido. La generación de ISO usa barra indeterminada porque no hay un porcentaje oficial que `eggs remaster` garantice.

### Instalación de una versión ya compilada

Si el directorio ya contiene el ejecutable `eggsmaker` generado previamente, en el equipo de destino no hace falta instalar Go ni compilar. Ejecutá:

```bash
chmod +x instalar-compilado.sh
./instalar-compilado.sh
```

El script instala el binario, el lanzador del menú, el icono y el README mediante `make install-files`. No ejecuta `go build`.
