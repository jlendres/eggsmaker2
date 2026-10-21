import asyncio
import getpass
import os
import re
import socket
import shutil
import tempfile
import time
import shlex
from nicegui import ui

from version import __app__, __version__
from backend import EggsmakerBackend

# Initialize backend
backend = EggsmakerBackend()

# Theme colors matching eggsmaker.py
COLORS = {
    'bg': '#23272e',          # Main background
    'panel': '#000000',       # Panels
    'button': '#0E48C5',      # Buttons
    'button_hover': '#1741a6',
    'orange': '#FD8637',      # Accents
    'terminal_bg': '#0e1010',
    'terminal_text': '#87CEFA',
    'light_blue': '#00BFFF',
    'success': '#39ee39',     # Green text
    'error': '#ff0000',
    'switch_on': '#00FF00',
    'switch_off': '#000000',
}

THEME_CSS = """
:root {{
  --egg-bg: {bg};
  --egg-panel: {panel};
  --egg-button: {button};
  --egg-button-hover: {button_hover};
  --egg-orange: {orange};
  --egg-terminal-bg: {terminal_bg};
  --egg-terminal-text: {terminal_text};
  --egg-light-blue: {light_blue};
  --egg-success: {success};
  --egg-error: {error};
  
  /* Button state colors */
  --btn-executing: #ff0000;
  --btn-success: #39ee39;
  --btn-disabled: #8b8b8b;
  --btn-warning: #FD8637;
}}

/* Global styles */
body {{ 
  background-color: var(--egg-bg) !important; 
  color: #ffffff; 
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}}

/* Smooth transitions for all elements */
* {{
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}}

/* Panel styling with glassmorphism */
.egg-panel {{ 
  background: linear-gradient(135deg, var(--egg-panel) 0%, rgba(0, 0, 0, 0.95) 100%) !important;
  border: 2px solid #333333;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
}}

/* Button base styling */
.egg-button {{ 
  background: linear-gradient(135deg, var(--egg-button) 0%, #0a3a8f 100%) !important;
  color: #ffffff !important; 
  font-weight: bold; 
  border-radius: 12px !important;
  box-shadow: 0 4px 12px rgba(14, 72, 197, 0.4);
  transition: all 0.3s ease;
  border: none;
}}

/* Global button rounding */
.q-btn {{
  border-radius: 12px !important;
}}

.egg-button:hover {{ 
  background: linear-gradient(135deg, var(--egg-button-hover) 0%, #0a2d6f 100%) !important;
  box-shadow: 0 6px 16px rgba(14, 72, 197, 0.6);
  transform: translateY(-2px);
}}

.egg-button:active {{
  transform: translateY(0px);
  box-shadow: 0 2px 8px rgba(14, 72, 197, 0.4);
}}

/* Button state: Executing (Red) */
.btn-executing {{
  background: linear-gradient(135deg, #ff0000 0%, #cc0000 100%) !important;
  box-shadow: 0 4px 12px rgba(255, 0, 0, 0.5);
  animation: pulse-red 2s infinite;
}}

/* Button state: Success (Green) */
.btn-success {{
  background: linear-gradient(135deg, #339991 0%, #2a7a73 100%) !important;
  box-shadow: 0 4px 12px rgba(51, 153, 145, 0.5);
}}

/* Button state: Disabled (Gray) */
.btn-disabled {{
  background: linear-gradient(135deg, #8b8b8b 0%, #6a6a6a 100%) !important;
  box-shadow: 0 2px 8px rgba(139, 139, 139, 0.3);
  opacity: 0.7;
  cursor: not-allowed !important;
}}

/* Pulse animation for executing state */
@keyframes pulse-red {{
  0%, 100% {{
    box-shadow: 0 4px 12px rgba(255, 0, 0, 0.5);
  }}
  50% {{
    box-shadow: 0 4px 20px rgba(255, 0, 0, 0.8);
  }}
}}

/* Title styling */
.egg-title {{ 
  color: var(--egg-light-blue) !important; 
  font-size: 17px; 
  font-weight: bold;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}}

/* Terminal styling */
.egg-terminal {{ 
  background-color: var(--egg-terminal-bg) !important; 
  color: var(--egg-terminal-text) !important; 
  font-family: 'JetBrains Mono', 'Consolas', 'Courier New', monospace;
  border-radius: 10px; 
  border: 1px solid #444C5E;
  box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.5);
  padding: 12px;
}}

/* Scrollbar styling for terminal */
.egg-terminal::-webkit-scrollbar {{
  width: 8px;
}}

.egg-terminal::-webkit-scrollbar-track {{
  background: rgba(0, 0, 0, 0.3);
  border-radius: 4px;
}}

.egg-terminal::-webkit-scrollbar-thumb {{
  background: rgba(135, 206, 250, 0.5);
  border-radius: 4px;
}}

.egg-terminal::-webkit-scrollbar-thumb:hover {{
  background: rgba(135, 206, 250, 0.8);
}}

/* Header styling */
.egg-header-title {{ 
  color: var(--egg-light-blue) !important; 
  font-weight: 800; 
  font-size: 1.4rem;
  text-shadow: 0 2px 8px rgba(0, 191, 255, 0.5);
}}

.q-page-container {{ padding-top: 0 !important; }}
.q-page {{ padding-top: 0 !important; }}
.q-header {{ 
  margin-bottom: 0 !important; 
  background: linear-gradient(135deg, var(--egg-panel) 0%, rgba(0, 0, 0, 0.98) 100%) !important;
  border-bottom: 1px solid rgba(68, 76, 94, 0.6);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}}

/* Progress bar styling */
.q-linear-progress {{
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid #444;
}}

.q-linear-progress__track {{
  background: rgba(68, 76, 94, 0.5) !important;
}}

/* Card hover effects */
.q-card {{
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}}

.q-card:hover {{
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}}

/* Switch styling */
.q-toggle {{
  transition: all 0.3s ease;
}}

.q-toggle__thumb {{
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}}

/* Input styling */
.q-field__control {{
  background: rgba(14, 16, 16, 0.6) !important;
  border-radius: 6px;
}}

.q-field__control:hover {{
  background: rgba(14, 16, 16, 0.8) !important;
}}

/* Dialog styling */
.q-dialog__backdrop {{
  backdrop-filter: blur(8px);
  background: rgba(0, 0, 0, 0.6) !important;
}}

.q-card {{
  border-radius: 12px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
}}

/* Animations */
@keyframes fadeIn {{
  from {{ opacity: 0; }}
  to {{ opacity: 1; }}
}}

@keyframes slideUp {{
  from {{ 
    opacity: 0;
    transform: translateY(20px);
  }}
  to {{ 
    opacity: 1;
    transform: translateY(0);
  }}
}}

.fade-in {{
  animation: fadeIn 0.5s ease;
}}

.slide-up {{
  animation: slideUp 0.5s ease;
}}

/* Status labels with glow effect */
.status-executing {{
  color: #ff0000 !important;
  text-shadow: 0 0 10px rgba(255, 0, 0, 0.5);
  animation: pulse-text 2s infinite;
}}

.status-success {{
  color: #39ee39 !important;
  text-shadow: 0 0 10px rgba(57, 238, 57, 0.5);
}}

@keyframes pulse-text {{
  0%, 100% {{ opacity: 1; }}
  50% {{ opacity: 0.7; }}
}}

/* Responsive improvements */
@media (max-width: 768px) {{
  .egg-header-title {{
    font-size: 1.1rem;
  }}
  .egg-panel {{
    border-radius: 8px;
  }}
}}
""".format(**COLORS)

# Global state
iso_generating = False
copying = False
copy_count = 0

copy_elapsed = 0
iso_elapsed = 0
copy_started_at = None
iso_started_at = None

# UI Elements
console = None
progress = None
phase_status = None
iso_size_label = None
copy_percent = None
copy_counter_label = None
copy_time_label = None
iso_time_label = None
Total_time_label = None

# UI Elements
versions_eggs = None
package_group_select = None
package_select = None
package_status = None
ui_states: dict[str, dict] = {}

PACKAGE_GROUPS = ('Nativos', 'AUR', 'Chaotic-AUR', 'Flatpak')

def current_ui_state() -> dict:
    try:
        return ui_states.get(ui.context.client.id, {})
    except RuntimeError:
        return {}

def append_log(text: str) -> None:
    client_console = current_ui_state().get('console')
    if client_console:
        client_console.push(text)

def set_progress(value: float | None, spinning: bool = False) -> None:
    client_progress = current_ui_state().get('progress')
    if client_progress is None:
        return
    if spinning:
        client_progress.props('indeterminate')
        client_progress.value = 0
    else:
        client_progress.props(remove='indeterminate')
        client_progress.value = max(0.0, min(1.0, value or 0.0))
    client_progress.update()

async def update_eggs_version() -> None:
    eggs_ver = 'No disponible'
    try:
        rc, out = await backend.run_capture('eggs --version')
        if rc == 0 and out:
            import re
            m = re.search(r"v?(\d+\.\d+(?:\.\d+)?)", out)
            if m:
                eggs_ver = m.group(1)
            else:
                eggs_ver = out.strip().splitlines()[-1]
    except Exception:
        pass

    label = current_ui_state().get('versions_eggs')
    if label:
        label.text = f"Penguins' Eggs: {eggs_ver}"
        label.update()

async def update_iso_size() -> None:
    iso_size_label = current_ui_state().get('iso_size_label')
    if not iso_size_label:
        return
    found_path = await find_latest_iso()
    if not found_path:
        iso_size_label.text = 'Tamaño ISO: No encontrada'
    else:
        try:
            rc, output = await backend.run_capture(f'stat -c %s -- {shlex.quote(found_path)}')
            if rc != 0:
                raise OSError(output.strip() or 'No se pudo consultar el tamaño de la ISO')
            iso_size_label.text = f"Tamaño ISO: {backend.format_size(int(output.strip()))}"
        except Exception as e:
            iso_size_label.text = f"Error: {e}"
    iso_size_label.update()

def check_sudo() -> bool:
    if not backend.sudo_available():
        ui.notify('No se encontró sudo. Instálelo para ejecutar operaciones administrativas.', type='negative')
        return False
    if not backend.sudo_password:
        ui.notify('Se requiere contraseña de sudo. Recargue la página o reinicie.', type='negative')
        return False
    return True

async def ask_yes_no(title: str, message: str) -> bool:
    future = asyncio.get_running_loop().create_future()
    with ui.dialog() as dialog:
        dialog.props('persistent')
        with ui.card().classes('egg-panel w-[520px] max-w-full'):
            ui.label(title).classes('text-xl font-bold text-white mb-2')
            ui.label(message).classes('text-white')

            def answer(value: bool) -> None:
                dialog.close()
                if not future.done():
                    future.set_result(value)

            with ui.row().classes('w-full justify-end mt-4'):
                ui.button('No', on_click=lambda: answer(False)).props('flat')
                ui.button('Sí', on_click=lambda: answer(True)).classes('egg-button')
    dialog.open()
    return await future

async def choose_copy_destination() -> str | None:
    future = asyncio.get_running_loop().create_future()
    with ui.dialog() as dialog:
        dialog.props('persistent')
        with ui.card().classes('egg-panel w-[650px] max-w-full'):
            ui.label('Destino de la ISO').classes('text-xl font-bold text-white mb-3')
            with ui.row().classes('w-full items-center'):
                destination = ui.input(value=backend.get_home_dir(), label='Carpeta').props('dark').classes('grow')

                async def browse_destination() -> None:
                    await open_dir_picker(destination)

                ui.button(
                    icon='folder',
                    on_click=browse_destination,
                ).props('flat round').tooltip('Examinar carpetas')

            def finish(value: str | None) -> None:
                dialog.close()
                if not future.done():
                    future.set_result(value)

            async def confirm() -> None:
                if not destination.value or not await backend.is_directory(destination.value):
                    ui.notify('La carpeta de destino no existe', type='warning')
                    return
                finish(destination.value)

            with ui.row().classes('w-full justify-end mt-4'):
                ui.button('Cancelar', on_click=lambda: finish(None)).props('flat')
                ui.button('Seleccionar', on_click=confirm).classes('egg-button')
    dialog.open()
    return await future

async def find_latest_iso() -> str | None:
    command = "find /home/eggs -maxdepth 2 -type f -name '*.iso' -printf '%T@ %p\\n' 2>/dev/null | sort -nr | head -n 1"
    rc, output = await backend.run_capture(command)
    if rc != 0 or not output.strip():
        return None
    return output.strip().split(' ', 1)[1]

async def copy_iso_to_destination(destination: str, iso_path: str | None) -> int:
    state = current_ui_state()
    phase_status = state.get('phase_status')
    progress = state.get('progress')
    copy_percent = state.get('copy_percent')
    copy_counter_label = state.get('copy_counter_label')
    if not await backend.command_exists('rsync'):
        ui.notify('Se requiere rsync para copiar la ISO', type='negative')
        return 127

    real_user = shlex.quote(backend.get_real_user())
    real_group = shlex.quote(backend.get_real_group())
    destination_arg = shlex.quote(destination.rstrip('/') + '/')
    ownership = f'--chown={real_user}:{real_group}'
    if iso_path:
        source = shlex.quote(iso_path)
        command = f'rsync {ownership} --info=progress2 -h {source} {destination_arg}'
    else:
        command = f'rsync {ownership} --info=progress2 -ah /home/eggs/ {destination_arg}'

    state['copying'] = True
    state['copy_started_at'] = time.monotonic()
    state.setdefault('copy_elapsed', 0.0)
    if progress:
        progress.props('color=deep-orange')
        progress.update()
    set_progress(0)
    phase_status.text = 'Ejecutando: Copiar ISO'
    phase_status.update()

    def log_copy_line(line: str) -> None:
        append_log(line)
        match = re.search(r'\b(\d{1,3})%', line)
        if match:
            percentage = min(100, int(match.group(1)))
            set_progress(percentage / 100)
            copy_percent.text = f'{percentage}%'
            copy_percent.update()

    append_log(f'$ {command}')
    try:
        rc = await backend.run_stream(command, log_copy_line)
        if rc == 0:
            state['copy_count'] = state.get('copy_count', 0) + 1
            copy_counter_label.text = f"Copias realizadas: {state['copy_count']}"
            copy_counter_label.update()
            ui.notify('Copia completada', type='positive')
        else:
            ui.notify(f'Error al copiar (código {rc})', type='negative')
        return rc
    finally:
        state['copying'] = False
        if state.get('copy_started_at') is not None:
            state['copy_elapsed'] = state.get('copy_elapsed', 0.0) + time.monotonic() - state['copy_started_at']
            state['copy_started_at'] = None
        set_progress(0)
        if progress:
            progress.props('color=light-blue')
            progress.update()
        copy_percent.text = '0%'
        copy_percent.update()
        phase_status.text = 'Estado: Esperando'
        phase_status.update()

async def run_iso_workflow() -> None:
    state = current_ui_state()
    phase_status = state.get('phase_status')
    if not check_sudo():
        return

    eggs = await backend.get_cmd_eggs()
    phase_status.text = 'Configurando creación de ISO'
    phase_status.update()

    if await ask_yes_no('Configuración del escritorio', '¿Desea copiar la configuración del escritorio (eggs tools skel)?'):
        set_progress(None, spinning=True)
        rc = await backend.run_stream(f'sudo {eggs} tools skel', append_log)
        set_progress(0)
        if rc != 0:
            phase_status.text = 'Estado: Esperando'
            phase_status.update()
            ui.notify('Falló la copia de configuración del escritorio', type='negative')
            return

    if await ask_yes_no('Configuración de Eggs', '¿Desea editar la configuración de Eggs?'):
        await edit_config_dialog()

    clone_system = await ask_yes_no('Tipo de ISO', '¿Desea crear un clon del sistema (--clone)?')
    if not await ask_yes_no('Confirmar generación', '¿Está seguro de iniciar la generación de la ISO?'):
        phase_status.text = 'Estado: Esperando'
        phase_status.update()
        return

    command = f'sudo {eggs} kill && {eggs} tools clean && {eggs} remaster'
    if clone_system:
        command += ' --clone'

    append_log('=== INICIANDO GENERACIÓN DE ISO ===')
    phase_status.text = 'Ejecutando: Generación de ISO'
    phase_status.update()
    state['iso_generating'] = True
    state['iso_started_at'] = time.monotonic()
    state.setdefault('iso_elapsed', 0.0)
    set_progress(None, spinning=True)
    try:
        rc = await backend.run_stream(command, append_log)
    finally:
        state['iso_generating'] = False
        if state.get('iso_started_at') is not None:
            state['iso_elapsed'] = state.get('iso_elapsed', 0.0) + time.monotonic() - state['iso_started_at']
            state['iso_started_at'] = None
        set_progress(0)
    if rc != 0:
        phase_status.text = 'Estado: Esperando'
        phase_status.update()
        ui.notify('Error al generar la ISO', type='negative')
        return

    iso_path = await find_latest_iso()
    await update_iso_size()
    if not iso_path:
        append_log('No se encontró una ISO en /home/eggs; se ofrecerá copiar el contenido de la carpeta.')

    copy_failed = False
    while True:
        destination = await choose_copy_destination()
        if destination is None:
            if await ask_yes_no('Omitir copia', '¿Desea omitir la copia de la ISO?'):
                break
            continue
        if await copy_iso_to_destination(destination, iso_path) != 0:
            copy_failed = True
            break
        if not await ask_yes_no('Copiar otra vez', '¿Desea copiar la ISO a otro destino?'):
            break

    cleanup_failed = False
    if await ask_yes_no('Limpiar archivos temporales', '¿Desea eliminar los archivos temporales de construcción (eggs kill)?'):
        phase_status.text = 'Ejecutando: Limpiar archivos temporales'
        phase_status.update()
        set_progress(None, spinning=True)
        cleanup_rc = await backend.run_stream(f'sudo {eggs} kill', append_log)
        cleanup_failed = cleanup_rc != 0
        set_progress(0)

    phase_status.text = 'Estado: Esperando'
    phase_status.update()
    if copy_failed or cleanup_failed:
        ui.notify('El proceso terminó con errores; revise el terminal', type='negative')
    else:
        ui.notify('Proceso finalizado', type='positive')

async def edit_config_dialog() -> bool:
    append_log("DEBUG: edit_config_dialog called")
    try:
        config_file = "/etc/penguins-eggs.d/custom.yaml"
        append_log(f"DEBUG: Reading {config_file}")
        rc, content = await backend.run_capture(f"sudo cat {shlex.quote(config_file)}")
        append_log(f"DEBUG: read rc={rc}, content_len={len(content) if content else 0}")
        if rc != 0 or not content:
            ui.notify('No se pudo leer custom.yaml. Verifique la autenticación sudo y el archivo del sistema.', type='warning')
            return False
    except Exception as e:
        append_log(f"DEBUG: edit_config_dialog error: {e}")
        ui.notify(f"Error: {e}", type='negative')
        return False
        
    try:
        future = asyncio.get_running_loop().create_future()
        with ui.dialog() as d:
            d.props('persistent')
            with ui.card().classes('w-[700px] max-w-full egg-panel'):
                ui.label('Editar configuración de Eggs').classes('text-xl font-bold mb-4')
                ui.label(config_file).classes('text-gray-400 text-sm')
                config_yaml = ui.textarea(value=content).props('dark autogrow input-class="text-white"').classes('w-full')
                
                async def on_save():
                    temp_path = None
                    try:
                        with tempfile.NamedTemporaryFile(
                            mode='w', encoding='utf-8', prefix='eggsmaker2-custom-',
                            dir=backend.get_home_dir(), delete=False,
                        ) as config_file_tmp:
                            temp_path = config_file_tmp.name
                            config_file_tmp.write(config_yaml.value.rstrip() + '\n')
                        rc = await backend.run_stream(f"sudo install -m 644 {shlex.quote(temp_path)} {shlex.quote(config_file)}", append_log)
                        if rc == 0:
                            ui.notify('Configuración guardada', type='positive')
                            d.close()
                            if not future.done():
                                future.set_result(True)
                        else:
                            ui.notify('No se pudo guardar la configuración', type='negative')
                    except Exception as e:
                        ui.notify(f'Error: {e}', type='negative')
                    finally:
                        if temp_path and os.path.exists(temp_path):
                            os.unlink(temp_path)

                def cancel_config() -> None:
                    d.close()
                    if not future.done():
                        future.set_result(False)
                        
                with ui.row().classes('w-full justify-end mt-4'):
                    ui.button('Cancelar', on_click=cancel_config).classes('mr-2 bg-red-600 hover:bg-red-700')
                    ui.button('Guardar cambios', on_click=lambda: asyncio.create_task(on_save())).classes('egg-button')
        append_log("DEBUG: Opening dialog")
        d.open()
        return await future
    except Exception as e:
        append_log(f"DEBUG: edit_config_dialog error: {e}")
        ui.notify(f"Error: {e}", type='negative')
        return False

async def run_live_installer(installer: str) -> None:
    if not check_sudo():
        return

    if installer == 'Calamares':
        command = 'sudo eggs sysinstall calamares'
        append_log('$ sudo eggs sysinstall calamares')
        phase_status.text = 'Ejecutando: Instalación Live con Calamares'
        phase_status.update()
        set_progress(None, spinning=True)
        try:
            rc = await backend.run_stream(command, append_log)
            ui.notify('Calamares finalizó' if rc == 0 else 'Calamares terminó con errores',
                      type='positive' if rc == 0 else 'negative')
        finally:
            set_progress(0)
            phase_status.text = 'Estado: Esperando'
            phase_status.update()
        return

    terminals = ('x-terminal-emulator', 'gnome-terminal', 'xfce4-terminal', 'konsole', 'xterm')
    terminal = next((name for name in terminals if shutil.which(name)), None)
    if not terminal:
        ui.notify('No se encontró una terminal gráfica para iniciar Krill', type='negative')
        return

    command = "sudo eggs sysinstall krill; printf '\\nKrill finalizó. '; read -r -p 'Presione Enter para cerrar...'"
    if terminal in ('gnome-terminal', 'x-terminal-emulator'):
        args = [terminal, '--', 'bash', '-lc', command]
    elif terminal == 'xfce4-terminal':
        args = [terminal, '--title=Eggsmaker 2 - Krill', '-e', f"bash -lc {shlex.quote(command)}"]
    elif terminal == 'konsole':
        args = [terminal, '--title', 'Eggsmaker 2 - Krill', '-e', 'bash', '-lc', command]
    else:
        args = [terminal, '-T', 'Eggsmaker 2 - Krill', '-e', 'bash', '-lc', command]

    try:
        process = await asyncio.create_subprocess_exec(*args)
        phase_status.text = 'Ejecutando: Instalación Live con Krill'
        phase_status.update()
        set_progress(None, spinning=True)
        append_log('$ sudo eggs sysinstall krill (terminal interactiva)')
        ui.notify('Krill se abrió en una terminal gráfica', type='positive')

        async def wait_for_krill() -> None:
            rc = await process.wait()
            set_progress(0)
            phase_status.text = 'Estado: Esperando'
            phase_status.update()
            if rc != 0:
                ui.notify('La terminal de Krill terminó con errores', type='negative')

        asyncio.create_task(wait_for_krill())
    except OSError as error:
        append_log(f'No se pudo abrir Krill: {error}')
        ui.notify(f'No se pudo abrir Krill: {error}', type='negative')

def package_list_dir() -> str:
    return os.path.join(backend.get_home_dir(), f'{socket.gethostname()}-pak')

def load_package_group(group: str) -> list[str]:
    filename = {
        'Nativos': 'lista-nativo.txt',
        'AUR': 'lista-aur.txt',
        'Chaotic-AUR': 'lista-chaotic-aur.txt',
        'Flatpak': 'lista-flatpak.txt',
    }.get(group)
    if not filename:
        return []
    try:
        with open(os.path.join(package_list_dir(), filename), encoding='utf-8') as package_file:
            return sorted({line.strip() for line in package_file if line.strip()})
    except OSError:
        return []

def refresh_package_options(state: dict | None = None) -> None:
    state = state or current_ui_state()
    package_group_select = state.get('package_group_select')
    package_select = state.get('package_select')
    package_status = state.get('package_status')
    if package_group_select is None or package_select is None:
        return
    packages = load_package_group(package_group_select.value or 'Nativos')
    package_select.options = packages
    package_select.value = []
    package_select.update()
    if package_status:
        package_status.text = f'{len(packages)} paquetes disponibles'
        package_status.update()

def clear_package_lists() -> None:
    state = current_ui_state()
    package_status = state.get('package_status')
    with ui.dialog() as dialog:
        with ui.card().classes('egg-panel'):
            ui.label('¿Vaciar las listas de paquetes recopiladas?').classes('text-white font-bold')
            ui.label(package_list_dir()).classes('text-gray-300')

            async def confirm_clear() -> None:
                list_files = ' '.join(
                    shlex.quote(os.path.join(package_list_dir(), filename))
                    for filename in (
                        'lista-nativo.txt', 'lista-aur.txt', 'lista-chaotic-aur.txt', 'lista-flatpak.txt'
                    )
                )
                rc = await backend.run_stream(f'rm -f {list_files}', append_log)
                dialog.close()
                refresh_package_options(state)
                package_status.text = 'Listas eliminadas' if rc == 0 else 'No se pudieron eliminar las listas'
                package_status.update()
                ui.notify(package_status.text, type='positive' if rc == 0 else 'negative')

            with ui.row().classes('w-full justify-end'):
                ui.button('Cancelar', on_click=dialog.close)
                ui.button('Vaciar listas', on_click=confirm_clear).props('color=negative')
    dialog.open()

async def collect_package_lists() -> None:
    state = current_ui_state()
    package_status = state.get('package_status')
    phase_status = state.get('phase_status')
    if not check_sudo():
        return

    package_status.text = 'Recopilando paquetes instalados...'
    package_status.update()
    phase_status.text = 'Ejecutando: Recopilar paquetes'
    phase_status.update()
    set_progress(None, spinning=True)
    try:
        _, release = await backend.run_capture('cat /etc/os-release')
        release_values = {}
        for line in release.splitlines():
            if '=' in line:
                key, value = line.split('=', 1)
                release_values[key] = value.strip().strip('"')
        distro_ids = f"{release_values.get('ID', '')} {release_values.get('ID_LIKE', '')}".lower()
        is_arch = any(name in distro_ids for name in ('arch', 'manjaro', 'endeavour', 'cachy'))
        is_debian = any(name in distro_ids for name in ('debian', 'ubuntu', 'mint', 'pop'))
        if not is_arch and not is_debian:
            raise RuntimeError('El gestor de paquetes solo admite familias Arch y Debian/Ubuntu')

        groups = {name: [] for name in PACKAGE_GROUPS}
        if is_arch:
            _, installed = await backend.run_capture("pacman -Qqe | sort | grep -vwE '^(base|base-devel)$'")
            _, aur = await backend.run_capture('pacman -Qqem | sort')
            _, chaotic_repo = await backend.run_capture(
                "pacman -Sl chaotic-aur 2>/dev/null | awk '$1 == \"chaotic-aur\" {print $2}'"
            )
            installed_names = set(installed.splitlines())
            aur_names = set(aur.splitlines())
            chaotic_names = set(chaotic_repo.splitlines())
            foreign_names = installed_names - aur_names
            groups['AUR'] = sorted(aur_names - chaotic_names)
            groups['Chaotic-AUR'] = sorted(foreign_names & chaotic_names)
            groups['Nativos'] = sorted(foreign_names - chaotic_names)
        else:
            _, native = await backend.run_capture('apt-mark showmanual | sort')
            groups['Nativos'] = sorted(set(native.splitlines()))

        user = shlex.quote(os.environ.get('SUDO_USER') or getpass.getuser())
        _, flatpak = await backend.run_capture(
            f"command sudo -u {user} -H env DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/$(id -u {user})/bus "
            f"XDG_RUNTIME_DIR=/run/user/$(id -u {user}) flatpak list --app --columns=application 2>/dev/null"
        )
        groups['Flatpak'] = sorted(set(flatpak.splitlines()))

        list_dir = package_list_dir()
        real_user = shlex.quote(backend.get_real_user())
        real_group = shlex.quote(backend.get_real_group())
        quoted_list_dir = shlex.quote(list_dir)
        list_files = ' '.join(
            shlex.quote(os.path.join(list_dir, filename))
            for filename in ('lista-nativo.txt', 'lista-aur.txt', 'lista-chaotic-aur.txt', 'lista-flatpak.txt')
        )
        rc = await backend.run_stream(
            f'install -d -o {real_user} -g {real_group} {quoted_list_dir} && '
            f'(chown {real_user}:{real_group} {list_files} 2>/dev/null || true)',
            append_log,
        )
        if rc != 0:
            raise PermissionError(f'No se pudo preparar la carpeta {list_dir}')
        filenames = {
            'Nativos': 'lista-nativo.txt',
            'AUR': 'lista-aur.txt',
            'Chaotic-AUR': 'lista-chaotic-aur.txt',
            'Flatpak': 'lista-flatpak.txt',
        }
        for group, filename in filenames.items():
            with open(os.path.join(list_dir, filename), 'w', encoding='utf-8') as package_file:
                package_file.write(''.join(f'{package}\n' for package in groups[group]))

        refresh_package_options(state)
        summary = ' · '.join(f'{group}: {len(packages)}' for group, packages in groups.items())
        package_status.text = f'Listas guardadas en {package_list_dir()} | {summary}'
        package_status.update()
        append_log(f'Listas de paquetes generadas: {summary}')
        ui.notify('Listas de paquetes actualizadas', type='positive')
    except Exception as error:
        append_log(f'Error al recopilar paquetes: {error}')
        package_status.text = f'Error: {error}'
        package_status.update()
        ui.notify(f'Error al recopilar paquetes: {error}', type='negative')
    finally:
        set_progress(0)
        phase_status.text = 'Estado: Esperando'
        phase_status.update()

async def manage_packages(remove: bool = False) -> None:
    state = current_ui_state()
    package_group_select = state.get('package_group_select')
    package_select = state.get('package_select')
    package_status = state.get('package_status')
    phase_status = state.get('phase_status')
    if not check_sudo():
        return
    group = package_group_select.value
    selected = package_select.value or []
    if not selected:
        ui.notify('Seleccione al menos un paquete', type='warning')
        return

    names = ' '.join(shlex.quote(package) for package in selected)
    user = shlex.quote(os.environ.get('SUDO_USER') or getpass.getuser())
    if group == 'Flatpak':
        action = 'uninstall -y --user' if remove else 'install -y --user'
        command = (
            f'command sudo -u {user} -H env DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/$(id -u {user})/bus '
            f'XDG_RUNTIME_DIR=/run/user/$(id -u {user}) flatpak {action} {names}'
        )
    elif group == 'AUR':
        if not await backend.command_exists('yay'):
            ui.notify('Se requiere yay para instalar o quitar paquetes AUR', type='warning')
            return
        action = '-Rns --noconfirm' if remove else '-S --needed --noconfirm'
        command = f'command sudo -u {user} -H yay {action} {names}'
    elif group == 'Chaotic-AUR':
        action = '-Rns --noconfirm' if remove else '-S --needed --noconfirm'
        package_names = selected if remove else [f'chaotic-aur/{package}' for package in selected]
        command = f"sudo pacman {action} " + ' '.join(shlex.quote(package) for package in package_names)
    else:
        action = '-Rns --noconfirm' if remove else '-S --needed --noconfirm'
        has_pacman = await backend.command_exists('pacman')
        command = f'sudo pacman {action} {names}' if has_pacman else (
            f'sudo apt-get remove -y {names}' if remove else f'sudo apt-get install -y {names}'
        )

    operation = 'Desinstalando' if remove else 'Instalando'
    phase_status.text = f'Ejecutando: {operation} paquetes ({group})'
    phase_status.update()
    package_status.text = f'{operation} {len(selected)} paquetes...'
    package_status.update()
    set_progress(None, spinning=True)
    append_log(f'$ {command}')
    try:
        rc = await backend.run_stream(command, append_log)
        package_status.text = 'Proceso completado' if rc == 0 else f'Proceso finalizado con código {rc}'
        package_status.update()
        ui.notify('Proceso completado' if rc == 0 else 'El gestor de paquetes informó errores',
                  type='positive' if rc == 0 else 'negative')
        if rc == 0:
            await collect_package_lists()
    finally:
        set_progress(0)
        phase_status.text = 'Estado: Esperando'
        phase_status.update()

async def open_dir_picker(target_input=None) -> None:
    append_log("DEBUG: open_dir_picker called")
    try:
        requested_path = target_input.value if target_input else backend.get_home_dir()
        current_path = requested_path if await backend.is_directory(requested_path) else backend.get_home_dir()
        append_log(f"DEBUG: current_path={current_path}")
    
        with ui.dialog() as d:
            with ui.card().classes('w-[600px] h-[500px] egg-panel flex flex-col'):
                ui.label('Seleccionar directorio destino').classes('text-xl font-bold mb-2 text-white')
                path_label = ui.label(current_path).classes('mb-2 text-gray-400 text-sm')
                
                # List container
                list_container = ui.column().classes('w-full grow overflow-y-auto border border-gray-600 rounded p-2')
                
                async def load_dir(path):
                    nonlocal current_path
                    current_path = path
                    path_label.text = path
                    list_container.clear()
                    
                    items = await backend.list_directories(path)
                    with list_container:
                        if not items and path != '/':
                             # Fallback if empty but not root, allow going up
                             items = [{'name': '..', 'path': os.path.dirname(path), 'type': 'dir'}]

                        for item in items:
                            icon = 'folder' if item['type'] == 'dir' else 'description'

                            async def navigate(_, selected_path=item['path']) -> None:
                                await load_dir(selected_path)

                            btn = ui.button(item['name'], icon=icon, on_click=navigate)
                            btn.props('flat align=left no-caps').classes('w-full text-left')

                
                await load_dir(current_path)
                
                with ui.row().classes('w-full justify-end mt-4'):
                    ui.button('Cancelar', on_click=d.close).classes('mr-2 bg-red-600 hover:bg-red-700')
                    def select():
                        if target_input:
                            target_input.value = current_path
                        d.close()
                    ui.button('Seleccionar', on_click=select).classes('egg-button')
        d.open()
    except Exception as e:
        append_log(f"DEBUG: open_dir_picker error: {e}")
        ui.notify(f"Error al abrir selector: {e}", type='negative')

def format_elapsed(seconds: float) -> str:
    whole_seconds = max(0, int(seconds))
    hours, remainder = divmod(whole_seconds, 3600)
    minutes, secs = divmod(remainder, 60)
    return f'{hours:02d}:{minutes:02d}:{secs:02d}'

async def update_timers() -> None:
    state = current_ui_state()
    now = time.monotonic()
    copy_total = state.get('copy_elapsed', 0.0) + (
        now - state['copy_started_at'] if state.get('copying') and state.get('copy_started_at') else 0
    )
    iso_total = state.get('iso_elapsed', 0.0) + (
        now - state['iso_started_at'] if state.get('iso_generating') and state.get('iso_started_at') else 0
    )
    copy_time_label = state.get('copy_time_label')
    iso_time_label = state.get('iso_time_label')
    total_time_label = state.get('total_time_label')
    if copy_time_label:
        copy_time_label.text = f"Copia: {format_elapsed(copy_total)}"
        copy_time_label.update()
    if iso_time_label:
        iso_time_label.text = f"Generación: {format_elapsed(iso_total)}"
        iso_time_label.update()
    total = copy_total + iso_total
    if total_time_label:
        total_time_label.text = f"Total: {format_elapsed(total)}"
        total_time_label.update()

@ui.page('/')
def main_page():
    global versions_eggs
    global phase_status, copy_percent, copy_counter_label, copy_time_label, iso_time_label, Total_time_label
    global progress, console
    global iso_size_label

    ui.add_head_html("<style>{}</style>".format(THEME_CSS))

    with ui.header().classes('items-center justify-between egg-panel relative'):
        with ui.row().classes('items-center'):
            logo_path = os.path.join(backend.base_path, 'assets', 'assets', 'eggsmaker2.png')
            if not os.path.exists(logo_path):
                logo_path = os.path.join(backend.base_path, 'assets', 'eggsmaker2.png')
            if not os.path.exists(logo_path):
                logo_path = os.path.join(backend.base_path, 'eggsmaker2.png')
            
            if os.path.exists(logo_path):
                ui.image(logo_path).props('fit=contain').style('width:32px;height:32px;border-radius:6px;')
            ui.label(f"{__app__} - Versión {__version__}").classes('text-bold')
            
        with ui.row().classes('items-center absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2'):
             ui.label('Eggsmaker 2').classes('egg-header-title ml-2').style('font-weight: bold !important; color: #FF5722 !important')
             
        with ui.column().classes('items-end'):
            versions_eggs = ui.label("Penguins' Eggs: N/A")

    # Terminal area
    with ui.row().classes('w-full q-pa-none q-ma-none q-mt-none main-top-row'):
        console = ui.log().classes('w-full h-[300px] egg-terminal p-2 overflow-y-auto')

    with ui.card().classes('w-full egg-panel q-mt-sm'):
        with ui.row().classes('w-full items-center justify-between'):
            ui.label('Generador de ISO').classes('egg-title')
            ui.button('Crear ISO', icon='disc_full', on_click=run_iso_workflow).classes('egg-button')

    async def start_live_installer() -> None:
        await run_live_installer(installer_choice.value)

    async def install_selected_packages() -> None:
        await manage_packages(remove=False)

    async def uninstall_selected_packages() -> None:
        await manage_packages(remove=True)

    with ui.row().classes('w-full items-stretch q-gutter-sm q-pa-none q-ma-none'):
        with ui.card().classes('min-w-[300px] grow egg-panel'):
            ui.label('Instalación Live').classes('egg-title')
            ui.label('Ejecuta Calamares en modo gráfico o abre Krill en una terminal interactiva.').classes('text-gray-300')
            installer_choice = ui.radio(['Calamares', 'Krill'], value='Calamares').props('inline color=light-blue')
            ui.button('Iniciar instalador', on_click=start_live_installer).classes('egg-button mt-2')

        with ui.card().classes('min-w-[300px] grow egg-panel'):
            ui.label('Configuración de Penguins’ Eggs').classes('egg-title')
            ui.label('Edita /etc/penguins-eggs.d/custom.yaml sin abrir el asistente de texto.').classes('text-gray-300')
            ui.button('Configurar Eggs', on_click=edit_config_dialog).classes('egg-button mt-2')

    # Status Bar - Simplified for visibility
    with ui.row().classes('w-full'):
        with ui.card().classes('w-full').style('background-color: #000000; padding: 16px; border: 2px solid #444C5E;'):
            phase_status = ui.label("Estado: Esperando").style('color: white; font-size: 18px; font-weight: bold; margin-bottom: 8px;')
            
            # Progress bar with explicit styling
            progress = ui.linear_progress(value=0, show_value=False).style('height: 24px; background-color: #444C5E; margin-bottom: 16px;')
            progress.props('color=light-blue')
            
            # Single info row below progress bar: left, center, right
            with ui.row().classes('w-full justify-between items-center'):
                # Left: Percentage and Copies
                with ui.row().classes('items-center').style('gap: 16px;'):
                    copy_percent = ui.label("0%").style('color: white; font-size: 18px; font-weight: bold;')
                    copy_counter_label = ui.label("Copias realizadas: 0").style('color: white; font-size: 14px; font-weight: bold;')
                
                # Center: ISO Size
                iso_size_label = ui.label('Tamaño ISO: N/A').style('color: #39ee39; font-size: 16px; font-weight: bold;')
                
                # Right: Timers
                with ui.row().classes('items-center').style('gap: 16px;'):
                    copy_time_label = ui.label("Copia: 00:00:00").style('color: cyan; font-size: 14px; font-weight: bold;')
                    iso_time_label = ui.label("Generación: 00:00:00").style('color: red; font-size: 14px; font-weight: bold;')
                    Total_time_label = ui.label("Total: 00:00:00").style('color: lime; font-size: 14px; font-weight: bold;')

    with ui.card().classes('w-full egg-panel q-mt-sm'):
        ui.label('Gestor de paquetes').classes('egg-title')
        ui.label('Recopila listas de paquetes del sistema y selecciona qué instalar o desinstalar.').classes('text-gray-300')
        with ui.row().classes('w-full items-end'):
            package_group_select = ui.select(
                list(PACKAGE_GROUPS), value='Nativos', label='Origen'
            ).props('dark').classes('min-w-[180px]')
            package_group_select.on_value_change(lambda _: refresh_package_options())
            package_select = ui.select(
                [], multiple=True, with_input=True, label='Paquetes'
            ).props('use-chips dark').classes('grow')
            ui.button('Recopilar', icon='refresh', on_click=collect_package_lists).classes('egg-button')
        with ui.row().classes('items-center'):
            ui.button('Instalar selección', icon='download', on_click=install_selected_packages).classes('egg-button')
            ui.button('Desinstalar selección', icon='delete', on_click=uninstall_selected_packages).props('color=negative')
            ui.button('Vaciar listas', icon='delete_sweep', on_click=clear_package_lists).props('flat')
            package_status = ui.label('Aún no hay listas. Use “Recopilar” para generarlas.').classes('text-gray-300')
                
    # Network info footer
    def show_about_dialog():
        with ui.dialog() as d, ui.card().classes('egg-panel w-[600px]'):
            ui.label('Acerca de Eggsmaker').classes('text-xl font-bold mb-4 text-white')
            
            ui.label('Eggsmaker creado por Jorge Luis Endres (c) 2025 Argentina').classes('text-white font-bold')
            with ui.row().classes('items-center'):
                 ui.label('Drive:').classes('text-gray-400 mr-2')
                 ui.link('https://drive.google.com/drive/folders/1rJfY1JaCKxaJD8XeFhYhyXROKZmaMe3B?usp=sharing', 'https://drive.google.com/drive/folders/1rJfY1JaCKxaJD8XeFhYhyXROKZmaMe3B?usp=sharing').classes('text-blue-400')
            ui.label('Feedback: jorge.eggsmaker@gmail.com').classes('text-gray-400 mb-4')
            
            ui.separator().classes('bg-gray-600 mb-4')
            
            ui.label('Penguins Eggs creado por Piero Proietti Italia').classes('text-white font-bold')
            with ui.row().classes('items-center'):
                 ui.label('Web:').classes('text-gray-400 mr-2')
                 ui.link('https://penguins-eggs.net', 'https://penguins-eggs.net').classes('text-blue-400')
            
            with ui.row().classes('w-full justify-end mt-6'):
                ui.button('Cerrar', on_click=d.close).classes('egg-button')
        d.open()

    with ui.footer().classes('egg-panel q-py-sm'):
        with ui.row().classes('w-full items-center relative-position'):
            # Left: Info button
            ui.button(icon='info', on_click=show_about_dialog).props('flat round dense color=white').tooltip('Acerca de')
            
            # Center: Network info
            access_host = backend.get_local_ip() if bind_host in ('0.0.0.0', '::') else bind_host
            ui.label(f"Acceso local: http://{access_host}:{port}").classes('absolute-center').style('color: #33b5e5 !important; font-weight: bold;')

    # Sudo Dialog (Created upfront)
    sudo_dialog = ui.dialog()
    with sudo_dialog, ui.card().classes('egg-panel'):
        ui.label('Autenticación requerida (sudo)').classes('text-white font-bold')
        if not backend.sudo_available():
            ui.label('No se detectó sudo en el sistema; instálelo para ejecutar operaciones administrativas.').classes('text-orange-400')
        pwd = ui.input('Contraseña de sudo', password=True, password_toggle_button=True).props('autofocus dark input-class="text-white"')

        async def authenticate() -> None:
            if not backend.sudo_available():
                ui.notify('No se encontró sudo. Instálelo antes de usar Eggsmaker 2.', type='negative')
                return
            password = pwd.value or ''
            valid, err = await backend.validate_sudo_password(password)
            if not valid:
                backend.sudo_password = None
                pwd.value = ''
                ui.notify(f'Error de autenticación: {err}', type='negative')
                return
            backend.sudo_password = password
            sudo_dialog.close()
            append_log('Autenticación sudo validada.')

        pwd.on('keydown.enter', authenticate)

        with ui.row():
            ui.button('Aceptar', on_click=authenticate).classes('egg-button')
            ui.button('Cancelar', on_click=sudo_dialog.close).classes('bg-red-600 hover:bg-red-700')

        client_id = ui.context.client.id
        ui_states[client_id] = {
            'console': console,
            'progress': progress,
            'phase_status': phase_status,
            'iso_size_label': iso_size_label,
            'copy_percent': copy_percent,
            'copy_counter_label': copy_counter_label,
            'copy_time_label': copy_time_label,
            'iso_time_label': iso_time_label,
            'total_time_label': Total_time_label,
            'versions_eggs': versions_eggs,
            'package_group_select': package_group_select,
            'package_select': package_select,
            'package_status': package_status,
            'copying': False,
            'iso_generating': False,
            'copy_count': 0,
            'copy_elapsed': 0.0,
            'iso_elapsed': 0.0,
            'copy_started_at': None,
            'iso_started_at': None,
        }
        refresh_package_options(ui_states[client_id])

    ui.timer(1.0, update_timers)
    ui.timer(0.2, update_eggs_version, once=True)
    ui.timer(0.05, sudo_dialog.open, once=True)

def main():
    """Entry point for the application"""
    global port, bind_host
    port = int(os.getenv('PORT', '8080'))
    bind_host = os.getenv('EGGSMAKER_HOST', '127.0.0.1')
    ui.run(title=f"{__app__}", port=port, host=bind_host, reload=False)

if __name__ == "__main__":
    main()
