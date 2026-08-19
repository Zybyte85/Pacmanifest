# Pacmanifest

Pacmanifest provides a simple, declarative way to manage Pacman packages.  
Define the packages you want in a manifest file and let Pacmanifest handle the installation and synchronization for you.

## Features

- **Declarative package list** – keep a single manifest (`manifest.pkgs`) that describes the exact set of packages you need.
- **Easy manifest management** – running `pacmanifest add` will automatically add packages to the manifest.
- **Sync with system** – `pacmanifest sync` ensures the system matches the manifest, removing packages that are no longer listed.
- **AUR Support** – packages from the AUR are supported via `yay`.
- **Version pinning** – packages can be pinned to a specific version in the manifest, even older ones.
- **Uses pacman** – uses the same package manager as Arch Linux, so it integrates very well with the system.

## Quick start

```sh
# Clone the repository
git clone https://github.com/Zybyte85/Pacmanifest
cd Pacmanifest

# Build the binary
go build -o pacmanifest .

# Create manifest and add a package
./pacmanifest add git 

# Add packages to the manifest manually
echo "vim\nhtop" >> ~/.pacmanifest/manifest.pkgs

# Add an AUR package to the manifest
./pacmanifest add aur:spotify

# Add an older version of a package to the manifest
./pacmanifest add cowsay=3.8.2

# Install the packages defined in the manifest
./pacmanifest sync
```

## Commands

| Command | Description |
|---------|-------------|
| `add` | Add a package to the manifest. |
| `sync`    | Sync the system state with the manifest (install missing, remove extra). |

## Configuration

- **Manifest file** – default location: `~/.pacmanifest/manifest.pkgs`. One package name per line.

## Contributing

Feel free to open issues or submit pull requests.

Pacmanifest is still in alpha. Some things, such as uninstalling packages that are removed from the manifest, are not yet supported, but planned.
