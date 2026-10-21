from importlib.util import module_from_spec, spec_from_file_location
from pathlib import Path
from setuptools import setup, find_packages

version_spec = spec_from_file_location('eggsmaker2_version', Path(__file__).with_name('version.py'))
version_module = module_from_spec(version_spec)
version_spec.loader.exec_module(version_module)
__version__ = version_module.__version__

with open("requirements.txt") as f:
    requirements = f.read().splitlines()

setup(
    name="eggsmaker2",
    version=__version__,
    author="Jorge Luis Endres",
    description="Graphical interface for Penguins' Eggs ISO creator",
    long_description=open("Readme Use/GUIA_DE_USO.md").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/pieroproietti/penguins-eggs",
    packages=find_packages(),
    py_modules=["main", "backend", "version"],
    include_package_data=True,
    install_requires=requirements,
    data_files=[
        ('share/applications', ['eggsmaker2.desktop']),
        ('share/pixmaps', ['assets/assets/eggsmaker2.png']),
        ('share/eggsmaker2/assets/assets', ['assets/assets/eggsmaker2.png']),
    ],
    entry_points={
        "console_scripts": [
            "eggsmaker2=main:main",
        ],
    },
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: System Administrators",
        "Operating System :: POSIX :: Linux",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Topic :: System :: Installation/Setup",
    ],
    python_requires=">=3.8",
)
