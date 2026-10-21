import asyncio
import getpass
import os
import re
import socket
import shutil
import tempfile
import time
import shlex
import yaml
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
ui_states: dict[str, dict] = {}

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

    if await ask_yes_no('Configuración de Eggs (sudo eggs config)', '¿Desea realizar cambios en la configuración de eggs (sudo eggs config)?'):
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
    config_file = "/etc/penguins-eggs.d/custom.yaml"

    try:
        rc, content = await backend.run_capture(f"sudo cat {shlex.quote(config_file)}")
        if rc != 0 or not content:
            content = (
                "# custom.yaml - Custom configurations for penguins-eggs\n"
                "remaster:\n"
                "  compression:\n"
                "    algorithm: zstd\n"
                "    level: 3\n"
            )
    except Exception as e:
        append_log(f"DEBUG: Error al leer {config_file}: {e}")
        content = (
            "# custom.yaml - Custom configurations for penguins-eggs\n"
            "remaster:\n"
            "  compression:\n"
            "    algorithm: zstd\n"
            "    level: 3\n"
        )

    try:
        data = yaml.safe_load(content) or {}
        if not isinstance(data, dict):
            data = {}
    except Exception:
        data = {}

    remaster = data.get('remaster') if isinstance(data.get('remaster'), dict) else {}
    init_user = remaster.get('user', '') or ''
    init_password = remaster.get('password', '') or ''

    compression = remaster.get('compression') if isinstance(remaster.get('compression'), dict) else {}
    init_algo = compression.get('algorithm', 'zstd') or 'zstd'
    if init_algo not in ('zstd', 'xz', 'gzip', 'lz4'):
        init_algo = 'zstd'
    try:
        init_level = int(compression.get('level', 3) or 3)
    except (ValueError, TypeError):
        init_level = 3

    try:
        future = asyncio.get_running_loop().create_future()
        with ui.dialog() as d:
            d.props('persistent')
            with ui.card().classes('w-[760px] max-w-[95vw] egg-panel p-5 flex flex-col gap-3'):
                # Header
                with ui.row().classes('w-full items-center justify-between border-b border-gray-700 pb-2'):
                    with ui.row().classes('items-center gap-2'):
                        ui.icon('tune', size='sm').style('color: #00BFFF;')
                        with ui.column().classes('gap-0'):
                            ui.label("Configuración de Penguins’ Eggs (sudo eggs config)").classes('text-lg font-bold text-white')
                            ui.label("Personalice la imagen ISO o ejecute el asistente oficial de Eggs").classes('text-xs text-gray-400')
                    ui.badge('custom.yaml', color='blue-grey-8').props('outline')

                # Navigation tabs
                with ui.tabs().classes('w-full text-white').props('dense indicator-color=primary align=left') as tabs:
                    tab_visual = ui.tab('visual', label='Opciones Principales', icon='tune')
                    tab_terminal = ui.tab('terminal', label='Asistente Terminal', icon='terminal')
                    tab_yaml = ui.tab('yaml', label='Editor YAML', icon='code')

                with ui.tab_panels(tabs, value=tab_visual).classes('w-full bg-transparent text-white'):
                    # TAB 1: Visual Form
                    with ui.tab_panel(tab_visual).classes('p-2 flex flex-col gap-4'):
                        with ui.card().classes('w-full bg-[#161a22] p-4 rounded-lg border border-gray-700'):
                            ui.label('Sesión en Vivo (Modo Live)').classes('text-sm font-bold text-cyan-400 mb-1')
                            ui.label('Configure las credenciales por defecto para iniciar sesión en el entorno Live de la ISO. Deje vacío para usar los valores predeterminados.').classes('text-xs text-gray-400 mb-3')
                            with ui.row().classes('w-full gap-4'):
                                user_input = ui.input(
                                    label='Usuario Live',
                                    value=init_user,
                                    placeholder='live'
                                ).props('dark outlined dense').classes('grow')
                                pass_input = ui.input(
                                    label='Contraseña Live',
                                    value=init_password,
                                    password=True,
                                    password_toggle_button=True,
                                    placeholder='evolution'
                                ).props('dark outlined dense').classes('grow')

                        with ui.card().classes('w-full bg-[#161a22] p-4 rounded-lg border border-gray-700'):
                            ui.label('Compresión de la Imagen (SquashFS)').classes('text-sm font-bold text-cyan-400 mb-1')
                            ui.label('Seleccione el algoritmo y nivel de compresión con el que se generará el sistema de archivos de la ISO.').classes('text-xs text-gray-400 mb-3')

                            algo_options = {
                                'zstd': 'zstd (Recomendado: rápido y equilibrado)',
                                'xz': 'xz (Máxima compresión: tamaño mínimo, más lento)',
                                'gzip': 'gzip (Compatibilidad estándar)',
                                'lz4': 'lz4 (Ultra rápido: compresión ligera)'
                            }
                            with ui.row().classes('w-full gap-4 items-center'):
                                algo_select = ui.select(
                                    algo_options,
                                    value=init_algo,
                                    label='Algoritmo'
                                ).props('dark outlined dense').classes('grow')

                                level_input = ui.number(
                                    label='Nivel ZSTD (1 - 22)',
                                    value=init_level,
                                    min=1,
                                    max=22,
                                    step=1
                                ).props('dark outlined dense').classes('w-44')
                                level_input.bind_visibility_from(algo_select, 'value', value='zstd')

                            algo_help = ui.label('').classes('text-xs text-gray-400 mt-1')
                            def update_algo_help():
                                hints = {
                                    'zstd': 'ℹ️ ZSTD comprime velozmente con muy buena reducción de tamaño. Nivel 3 es el estándar; 19 ofrece máxima compresión.',
                                    'xz': 'ℹ️ XZ ofrece el menor tamaño de ISO posible, pero requiere más tiempo y CPU durante la generación.',
                                    'gzip': 'ℹ️ GZIP es el compresor tradicional, universalmente compatible con todos los sistemas.',
                                    'lz4': 'ℹ️ LZ4 prioriza el tiempo de descompresión y arranque inmediato, con un tamaño de ISO mayor.'
                                }
                                algo_help.text = hints.get(algo_select.value, '')
                            algo_select.on_value_change(update_algo_help)
                            update_algo_help()

                    # TAB 2: Terminal Assistant
                    with ui.tab_panel(tab_terminal).classes('p-2 flex flex-col gap-3'):
                        with ui.card().classes('w-full bg-[#161a22] p-4 rounded-lg border border-gray-700'):
                            ui.label('Asistente interactivo en terminal (sudo eggs config)').classes('text-sm font-bold text-cyan-400 mb-2')
                            ui.label(
                                'Penguins’ Eggs dispone de un asistente interactivo oficial en modo texto. '
                                'Al presionar el botón se abrirá una ventana de terminal con dicho asistente. '
                                'Al finalizar y cerrar la terminal, los cambios se recargarán automáticamente aquí.'
                            ).classes('text-xs text-gray-300 leading-relaxed mb-4')

                            terminal_status = ui.label('').classes('text-xs text-yellow-400')

                            async def run_terminal_config():
                                terminals = ('konsole', 'x-terminal-emulator', 'xfce4-terminal', 'gnome-terminal', 'xterm')
                                terminal = next((name for name in terminals if shutil.which(name)), None)
                                if not terminal:
                                    ui.notify('No se encontró un emulador de terminal gráfica en el sistema', type='negative')
                                    return

                                terminal_status.text = 'Terminal abierta ejecutando sudo eggs config... Al terminar vuelva aquí.'
                                cmd = "sudo eggs config; echo; read -r -p 'Presione Enter para volver a Eggsmaker...'"
                                if terminal in ('gnome-terminal', 'x-terminal-emulator'):
                                    args = [terminal, '--', 'bash', '-lc', cmd]
                                elif terminal == 'xfce4-terminal':
                                    args = [terminal, '--title=Eggsmaker 2 - eggs config', '-e', f"bash -lc {shlex.quote(cmd)}"]
                                elif terminal == 'konsole':
                                    args = [terminal, '--title', 'Eggsmaker 2 - eggs config', '-e', 'bash', '-lc', cmd]
                                else:
                                    args = [terminal, '-T', 'Eggsmaker 2 - eggs config', '-e', 'bash', '-lc', cmd]

                                try:
                                    proc = await asyncio.create_subprocess_exec(*args)
                                    await proc.wait()
                                    terminal_status.text = 'Asistente finalizado. Configuración recargada exitosamente.'
                                    ui.notify('Configuración actualizada desde eggs config', type='positive')

                                    _, new_raw = await backend.run_capture(f"sudo cat {shlex.quote(config_file)}")
                                    if new_raw:
                                        yaml_editor.value = new_raw
                                        try:
                                            new_d = yaml.safe_load(new_raw) or {}
                                            if isinstance(new_d, dict):
                                                rem = new_d.get('remaster') or {}
                                                user_input.value = rem.get('user', '') or ''
                                                pass_input.value = rem.get('password', '') or ''
                                                cmp_info = rem.get('compression') or {}
                                                if cmp_info.get('algorithm') in algo_options:
                                                    algo_select.value = cmp_info.get('algorithm')
                                                if 'level' in cmp_info:
                                                    try:
                                                        level_input.value = int(cmp_info.get('level'))
                                                    except (ValueError, TypeError):
                                                        pass
                                        except Exception:
                                            pass
                                except Exception as err:
                                    terminal_status.text = f'Error al iniciar terminal: {err}'
                                    ui.notify(f'Error al abrir terminal: {err}', type='negative')

                            with ui.row().classes('items-center gap-3'):
                                ui.button('Abrir sudo eggs config en Terminal', icon='open_in_new', on_click=lambda: asyncio.create_task(run_terminal_config())).classes('egg-button')

                    # TAB 3: Advanced YAML Editor
                    with ui.tab_panel(tab_yaml).classes('p-2 flex flex-col gap-2'):
                        ui.label('Edición directa de custom.yaml').classes('text-sm font-bold text-cyan-400')
                        ui.label('Edición manual avanzada para parámetros personalizados de /etc/penguins-eggs.d/custom.yaml.').classes('text-xs text-gray-400')
                        yaml_editor = ui.textarea(value=content).props('dark outlined input-class="font-mono text-xs text-cyan-100"').classes('w-full').style('height: 240px;')

                # Footer / Actions
                with ui.row().classes('w-full justify-between items-center border-t border-gray-700 pt-3'):
                    ui.label(config_file).classes('text-xs text-gray-500 font-mono')

                    with ui.row().classes('items-center gap-2'):
                        def cancel_config() -> None:
                            d.close()
                            if not future.done():
                                future.set_result(False)

                        async def on_save():
                            temp_path = None
                            try:
                                if tabs.value == 'visual':
                                    current_data = {}
                                    try:
                                        current_data = yaml.safe_load(yaml_editor.value or content) or {}
                                        if not isinstance(current_data, dict):
                                            current_data = {}
                                    except Exception:
                                        current_data = {}

                                    if 'remaster' not in current_data or not isinstance(current_data['remaster'], dict):
                                        current_data['remaster'] = {}

                                    if user_input.value and user_input.value.strip():
                                        current_data['remaster']['user'] = user_input.value.strip()
                                    elif 'user' in current_data['remaster']:
                                        del current_data['remaster']['user']

                                    if pass_input.value and pass_input.value.strip():
                                        current_data['remaster']['password'] = pass_input.value.strip()
                                    elif 'password' in current_data['remaster']:
                                        del current_data['remaster']['password']

                                    if 'compression' not in current_data['remaster'] or not isinstance(current_data['remaster']['compression'], dict):
                                        current_data['remaster']['compression'] = {}

                                    current_data['remaster']['compression']['algorithm'] = algo_select.value
                                    if algo_select.value == 'zstd':
                                        current_data['remaster']['compression']['level'] = int(level_input.value or 3)
                                    elif 'level' in current_data['remaster']['compression']:
                                        del current_data['remaster']['compression']['level']

                                    out_text = yaml.dump(current_data, sort_keys=False, default_flow_style=False)
                                else:
                                    out_text = yaml_editor.value.rstrip() + '\n'
                                    try:
                                        yaml.safe_load(out_text)
                                    except Exception as y_err:
                                        ui.notify(f'Error de sintaxis YAML: {y_err}', type='negative')
                                        return

                                with tempfile.NamedTemporaryFile(
                                    mode='w', encoding='utf-8', prefix='eggsmaker2-custom-',
                                    dir=backend.get_home_dir(), delete=False,
                                ) as config_file_tmp:
                                    temp_path = config_file_tmp.name
                                    config_file_tmp.write(out_text)

                                rc = await backend.run_stream(f"sudo install -m 644 {shlex.quote(temp_path)} {shlex.quote(config_file)}", append_log)
                                if rc == 0:
                                    ui.notify('Configuración guardada correctamente', type='positive')
                                    d.close()
                                    if not future.done():
                                        future.set_result(True)
                                else:
                                    ui.notify('No se pudo guardar la configuración', type='negative')
                            except Exception as e:
                                ui.notify(f'Error al guardar: {e}', type='negative')
                            finally:
                                if temp_path and os.path.exists(temp_path):
                                    os.unlink(temp_path)

                        ui.button('Cancelar', on_click=cancel_config).classes('bg-red-600 hover:bg-red-700 text-white')
                        ui.button('Guardar cambios', icon='save', on_click=lambda: asyncio.create_task(on_save())).classes('egg-button')

            d.open()
            return await future
    except Exception as e:
        append_log(f"DEBUG: edit_config_dialog error: {e}")
        ui.notify(f"Error: {e}", type='negative')
        return False

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
            'copying': False,
            'iso_generating': False,
            'copy_count': 0,
            'copy_elapsed': 0.0,
            'iso_elapsed': 0.0,
            'copy_started_at': None,
            'iso_started_at': None,
        }

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
