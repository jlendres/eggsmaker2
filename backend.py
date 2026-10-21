import asyncio
import getpass
import grp
import os
import pwd
import shlex
import shutil
import socket
import sys
import gettext
import locale
from typing import Tuple

# Import version info
try:
    from version import __version__, __app__
except ImportError:
    __version__ = "0.0.0"
    __app__ = "Eggsmaker"

class EggsmakerBackend:
    def __init__(self):
        self.sudo_password: str | None = None
        self.base_path = self.get_base_path()
        self._setup_i18n()

    @property
    def in_flatpak_sandbox(self) -> bool:
        return os.path.exists('/.flatpak-info') or bool(os.environ.get('FLATPAK_ID'))

    def sudo_available(self) -> bool:
        return bool(shutil.which('sudo') or (self.in_flatpak_sandbox and shutil.which('flatpak-spawn')))

    async def command_exists(self, command: str) -> bool:
        rc, output = await self.run_capture(f'command -v {shlex.quote(command)}')
        return rc == 0 and bool(output.strip())

    async def validate_sudo_password(self, password: str) -> Tuple[bool, str]:
        # Si sudo no requiere contraseña (NOPASSWD)
        try:
            check_nopass = ['sudo', '-n', 'true']
            if self.in_flatpak_sandbox:
                check_nopass = ['flatpak-spawn', '--host', *check_nopass]
            p_nopass = await asyncio.create_subprocess_exec(
                *check_nopass,
                stdout=asyncio.subprocess.DEVNULL,
                stderr=asyncio.subprocess.DEVNULL,
            )
            await p_nopass.wait()
            if p_nopass.returncode == 0:
                return True, ""
        except OSError:
            pass

        if not password:
            return False, "Por favor ingrese la contraseña de sudo"

        # Invalidar cualquier sesión o timestamp residual previo
        try:
            reset_cmd = ['sudo', '-k']
            if self.in_flatpak_sandbox:
                reset_cmd = ['flatpak-spawn', '--host', *reset_cmd]
            p_reset = await asyncio.create_subprocess_exec(
                *reset_cmd,
                stdout=asyncio.subprocess.DEVNULL,
                stderr=asyncio.subprocess.DEVNULL,
            )
            await p_reset.wait()
        except OSError:
            pass

        try:
            command = ['sudo', '-S', '-p', '', '-v']
            if self.in_flatpak_sandbox:
                command = ['flatpak-spawn', '--host', *command]
            process = await asyncio.create_subprocess_exec(
                *command,
                stdin=asyncio.subprocess.PIPE,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE,
            )
            _, stderr = await process.communicate(input=(password + '\n').encode())
            if process.returncode == 0:
                return True, ""
            err_msg = stderr.decode(errors='ignore').strip()
            return False, err_msg or "Contraseña incorrecta o permisos insuficientes"
        except OSError as e:
            return False, str(e)

    def get_base_path(self) -> str:
        """Obtiene la ruta base de la aplicación, funciona tanto en desarrollo como compilado"""
        try:
            # Si estamos en un entorno de PyInstaller
            if getattr(sys, 'frozen', False):
                if hasattr(sys, '_MEIPASS'):
                    meipass = sys._MEIPASS
                    # Verificar rutas comunes en _MEIPASS
                    if os.path.exists(os.path.join(meipass, 'assets', 'assets', 'eggsmaker2.png')):
                        return meipass
                    if os.path.exists(os.path.join(meipass, 'eggsmaker.png')):
                        return meipass
            
            # Lista de posibles ubicaciones
            possible_paths = [
                os.path.dirname(os.path.abspath(__file__)),
                os.getcwd(),
                "/usr/share/eggsmaker",
                "/usr/local/share/eggsmaker",
            ]

            for path in possible_paths:
                if os.path.exists(os.path.join(path, "assets", "assets", "eggsmaker2.png")):
                    return path
                if os.path.exists(os.path.join(path, "eggsmaker.png")):
                    return path
            
            return os.path.dirname(os.path.abspath(__file__))
        except Exception:
            return os.path.dirname(os.path.abspath(__file__))

    def _setup_i18n(self):
        # Placeholder for future i18n implementation
        pass

    def tr(self, text: str) -> str:
        """Translate text (placeholder)"""
        return text

    def get_local_ip(self) -> str:
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            # doesn't even have to be reachable
            s.connect(('10.255.255.255', 1))
            IP = s.getsockname()[0]
            s.close()
        except Exception:
            IP = '127.0.0.1'
        return IP

    async def run_stream(self, cmd: str, log_callback) -> int:
        wrapped = cmd
        stdin = None
        password_input = None
        if self.sudo_password:
            cmd_no_sudo = cmd.lstrip()
            if cmd_no_sudo.startswith('sudo '):
                cmd_no_sudo = cmd_no_sudo[len('sudo '):]
            safe_path = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
            # Sanitize environment variables to prevent AppImage interference
            env_vars = "LD_PRELOAD= LD_LIBRARY_PATH= APPIMAGE= APPDIR= PYTHONPATH="
            wrapped = (
                f"sudo -S -E -p '' env {env_vars} PATH={safe_path} bash -lc {shlex.quote(cmd_no_sudo)}"
            )
            stdin = asyncio.subprocess.PIPE
            password_input = (self.sudo_password + '\n').encode()

        if self.in_flatpak_sandbox:
            wrapped = f'flatpak-spawn --host /bin/bash -lc {shlex.quote(wrapped)}'
        
        if log_callback:
            log_callback(f"$ {'sudo ' if password_input else ''}{cmd_no_sudo if password_input else cmd}")
            
        process = await asyncio.create_subprocess_shell(
            wrapped,
            stdin=stdin,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.STDOUT,
            executable="/bin/bash",
        )

        if password_input and process.stdin:
            process.stdin.write(password_input)
            await process.stdin.drain()
            process.stdin.close()
        
        # Split on both newlines and carriage returns so progress tools can stream updates.
        if process.stdout:
            pending = bytearray()
            while True:
                chunk = await process.stdout.read(4096)
                if not chunk:
                    break
                pending.extend(chunk)
                while True:
                    delimiters = [index for index in (pending.find(b'\n'), pending.find(b'\r')) if index >= 0]
                    if not delimiters:
                        break
                    end = min(delimiters)
                    decoded_line = pending[:end].decode(errors='ignore').rstrip()
                    del pending[:end + 1]
                    if log_callback:
                        log_callback(decoded_line)
                    await asyncio.sleep(0)
            if pending and log_callback:
                log_callback(pending.decode(errors='ignore').rstrip())
                    
        rc = await process.wait()
        return rc

    async def run_capture(self, cmd: str) -> Tuple[int, str]:
        wrapped = cmd
        password_input = None
        if self.sudo_password:
            cmd_no_sudo = cmd.lstrip()
            if cmd_no_sudo.startswith('sudo '):
                cmd_no_sudo = cmd_no_sudo[len('sudo '):]
            safe_path = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
            # Sanitize environment variables to prevent AppImage interference
            env_vars = "LD_PRELOAD= LD_LIBRARY_PATH= APPIMAGE= APPDIR= PYTHONPATH="
            wrapped = (
                f"sudo -S -E -p '' env {env_vars} PATH={safe_path} bash -lc {shlex.quote(cmd_no_sudo)}"
            )
            password_input = (self.sudo_password + '\n').encode()
        if self.in_flatpak_sandbox:
            wrapped = f'flatpak-spawn --host /bin/bash -lc {shlex.quote(wrapped)}'
        try:
            process = await asyncio.create_subprocess_shell(
                wrapped,
                stdin=asyncio.subprocess.PIPE if password_input else None,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE,
                executable="/bin/bash",
            )
            stdout, stderr = await process.communicate(input=password_input)
            
            try:
                out_str = stdout.decode(errors='ignore')
            except:
                out_str = stdout.decode('utf-8', errors='ignore')

            if process.returncode != 0 and stderr:
                error_str = stderr.decode(errors='ignore').strip()
                if error_str:
                    out_str = f'{out_str.rstrip()}\n{error_str}'.strip()
            
            # print(f"DEBUG: run_capture cmd={cmd} rc={process.returncode}") 
            return process.returncode, out_str
        except Exception as e:
            print(f"DEBUG: run_capture exception: {e}")
            return -1, str(e)

    async def get_cmd_eggs(self) -> str:
        try:
            rc, output = await self.run_capture('command -v eggs')
            if rc == 0 and output.strip():
                return output.strip()
            return 'eggs'
        except Exception:
            return 'eggs'

    def format_size(self, size_bytes: int) -> str:
        if size_bytes < 1024:
            return f"{size_bytes} B"
        if size_bytes < 1024 * 1024:
            return f"{size_bytes / 1024:.2f} KB"
        if size_bytes < 1024 * 1024 * 1024:
            return f"{size_bytes / (1024 * 1024):.2f} MB"
        return f"{size_bytes / (1024 * 1024 * 1024):.2f} GB"

    def get_home_dir(self) -> str:
        try:
            return pwd.getpwnam(self.get_real_user()).pw_dir
        except KeyError:
            pass
        return os.path.expanduser('~')

    def get_real_user(self) -> str:
        return os.environ.get('SUDO_USER') or getpass.getuser()

    def get_real_user_ids(self) -> tuple[int, int]:
        try:
            account = pwd.getpwnam(self.get_real_user())
            return account.pw_uid, account.pw_gid
        except KeyError:
            return os.getuid(), os.getgid()

    def get_real_group(self) -> str:
        _, group_id = self.get_real_user_ids()
        return grp.getgrgid(group_id).gr_name

    async def is_directory(self, path: str) -> bool:
        rc, _ = await self.run_capture(f'test -d {shlex.quote(path)}')
        return rc == 0

    async def list_directories(self, path: str) -> list[dict]:
        normalized_path = os.path.normpath(path)
        command = (
            f"find {shlex.quote(normalized_path)} -mindepth 1 -maxdepth 1 -type d "
            "! -name '.*' -printf '%f\\0' 2>/dev/null"
        )
        rc, output = await self.run_capture(command)
        if rc != 0:
            raise PermissionError(f'No se pudo leer el directorio: {normalized_path}')

        names = (name for name in output.split('\0') if name)
        items = [
            {'name': name, 'path': os.path.join(normalized_path, name), 'type': 'dir'}
            for name in names
        ]
        items.sort(key=lambda item: item['name'].lower())
        if normalized_path != '/':
            items.insert(0, {
                'name': '..',
                'path': os.path.dirname(normalized_path),
                'type': 'dir',
            })
        return items
