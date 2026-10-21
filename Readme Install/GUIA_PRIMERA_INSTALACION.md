# Guía de Usuario: Primera Instalación y Portabilidad

Esta guía detalla los archivos necesarios para llevar este proyecto a otros equipos y cómo realizar la instalación.

## 1. Archivos Necesarios para Portar

Para instalar este proyecto en otros equipos, debes copiar la carpeta completa del proyecto. Asegúrate de que incluya los siguientes archivos y carpetas esenciales:

### Código Fuente y Recursos
*   **`main.py`**: Archivo principal de la aplicación.
*   **`backend.py`**: Lógica del sistema.
*   **`version.py`**: Versión de la aplicación.
*   **`assets/`**: Carpeta con imágenes y recursos (logo, etc.).

### Instalación
*   **`install.sh`**: Script de instalación automatizada.
*   **`requirements.txt`**: Dependencias de Python.
*   **`setup.py`** y **`pyproject.toml`**: Configuración del paquete.
*   **`eggsmaker2.desktop`**: Acceso directo para el menú.

---

## 2. Instrucciones de Instalación

Una vez copiada la carpeta del proyecto al nuevo equipo, tienes dos opciones de instalación dependiendo de tus necesidades.

### Opción A: Instalación Local (Recomendada)
Ideal para un solo usuario. No afecta al sistema operativo globalmente.

1.  Abre una terminal dentro de la carpeta del proyecto.
2.  Ejecuta el siguiente comando:
    ```bash
    ./install.sh
    ```

**Ubicación de archivos:**
*   El programa se instalará en: `~/.local/share/eggsmaker2/`
*   El ejecutable estará en: `~/.local/bin/eggsmaker2`

### Opción B: Instalación en el Sistema
Hace que el programa esté disponible para todos los usuarios del equipo. Requiere permisos de administrador.

1.  Abre una terminal dentro de la carpeta del proyecto.
2.  Ejecuta el siguiente comando con `sudo`:
    ```bash
    sudo ./install.sh --system
    ```

**Ubicación de archivos:**
*   El programa se instalará en: `/opt/eggsmaker2/`
*   El ejecutable estará en: `/usr/local/bin/eggsmaker2`

---

## 3. Ejecución

Una vez instalado, puedes buscar **"Eggsmaker 2"** en el menú de aplicaciones de tu sistema o ejecutarlo desde la terminal con el comando:

```bash
eggsmaker2
```

---

## 4. Actualización

Si necesitas actualizar la aplicación a una nueva versión:

1.  Descarga o copia la nueva versión del proyecto.
2.  Repite el paso de instalación (Opción A o B) que usaste originalmente.
    *   El instalador sobrescribirá automáticamente los archivos antiguos con los nuevos.
    *   No perderás tus configuraciones personales si están fuera de la carpeta del programa.

---

## 5. Desinstalación

Para eliminar la aplicación por completo:

1.  Abre una terminal en la carpeta del proyecto (si aún la tienes) o descarga el script `uninstall.sh`.
2.  Ejecuta el desinstalador según tu tipo de instalación:

**Si instalaste localmente (Opción A):**
```bash
./uninstall.sh
```

**Si instalaste en el sistema (Opción B):**
```bash
sudo ./uninstall.sh --system
```


