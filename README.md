## **Klir**
Interactive CLI to bulk stop & delete Docker containers and images.

### **Usage**
```bash
klir [global options] <command> [command options]
```

### **Global Options**
- `-h`, `--help`   Show help  
- `-v`, `--version` Print version


### **Commands**

#### `klir stop`
Interactively stop running containers.

```bash
klir stop [OPTIONS]
```
Options:
- `-v`, `--verbose` Show debug logs


#### `klir rm`
Interactively remove **inactive** containers.

```bash
klir rm [OPTIONS]
```
Options:
- `-f`, `--force`  Allow removal of **running** containers  
- `-v`, `--verbose` Show debug logs


#### `klir rmi`
Interactively remove **orphaned** images.

```bash
klir rmi [OPTIONS]
```
Options:
- `-f`, `--force`  Allow removal of **active** images with stopped containers 
- `-v`, `--verbose` Show debug logs

### Installation

#### Option 1: Download Binary
- Visit the [Releases](https://github.com/melsonic/klir/releases) page.
- Download the appropriate binary from the **Assets** section.
- Run the binary directly.
##### Linux / macOS (Terminal)
```bash
./klir

# Optional: make it available system-wide
sudo cp ./klir /usr/local/bin
```

##### Windows (Command Prompt)
```cmd
.\klir.exe
```

> **Note:** To run `klir.exe` system-wide on Windows, add its location to the [Environment Variables](https://medium.com/@kevinmarkvi/how-to-add-executables-to-your-path-in-windows-5ffa4ce61a53).

#### Option 2: Build from Source
- Ensure [Golang](https://go.dev/doc/install) & [Docker](https://docs.docker.com/engine/install/) is installed on your system.
- Run the following commands in order
```
    git clone https://github.com/melsonic/klir.git
    cd klir
    go build -o klir main.go docker_client.go
    ./klir
````

### Result

<img src="./assets/klir.gif" style="border: 1px solid white; border-radius: 5px;" alt="description">