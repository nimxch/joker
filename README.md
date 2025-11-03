# JokerDB 🃏

> A Redis-inspired in-memory key-value store implementation in Go

JokerDB is a lightweight, Redis-inspired key-value database implemented in Go. It provides a TCP-based server with a simple TLV (Tag-Length-Value) protocol for storing and retrieving data in memory.

## ✨ Features

- **🔑 Key-Value Store**: Simple and efficient in-memory key-value storage
- **⚡ TCP Server**: Network-accessible server listening on port 9999
- **📦 TLV Protocol**: Custom Tag-Length-Value protocol for client-server communication
- **🚀 Lightweight**: Minimal dependencies with pure Go implementation
- **🔧 Concurrent**: Handle multiple client connections simultaneously

## 📋 Prerequisites

- Go 1.20 or higher
- Make (optional, for using Makefile commands)

## 🚀 Getting Started

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/nimaidev/joker.git
   cd joker
   ```

2. **Build the project:**

   Using Make:
   ```bash
   make build
   ```

   Or using Go directly:
   ```bash
   go build -o .build/joker
   ```

3. **Run JokerDB:**

   Using Make:
   ```bash
   make run
   ```

   Or run the binary directly:
   ```bash
   ./.build/joker
   ```

   The server will start and listen on `localhost:9999`.

## 💡 Usage

JokerDB uses a custom TLV (Tag-Length-Value) protocol for communication. The protocol format is:

```
[Tag: 2 bytes][Length: 2 bytes][Value: N bytes]
```

### Supported Operations

#### PUT Operation (Tag: 1)
Store a key-value pair in the database.

**Format:** `key>value`

**Example:**
```
Tag: 0x0001 (PUT command)
Length: length of "mykey>myvalue"
Value: "mykey>myvalue"
```

#### GET Operation (Tag: 2)
Retrieve a value by its key.

**Format:** `key`

**Example:**
```
Tag: 0x0002 (GET command)
Length: length of "mykey"
Value: "mykey"
```

### Response Codes

- `0`: Success
- `101`: Key not found error

### Example Client Connection

You can connect to JokerDB using any TCP client. Here's an example using `nc` (netcat):

```bash
# Connect to the server
nc localhost 9999

# Then send binary data according to the TLV protocol
```

For a more practical approach, you'll need to implement a client that properly encodes data using the TLV protocol.

## 🏗️ Project Structure

```
.
├── main.go           # Main entry point and TCP server
├── server/           # Server implementation with TLV protocol
├── parser/           # Command parsing logic
├── constants/        # Application constants
├── utils/            # Utility functions
├── config/           # Configuration files
├── Makefile          # Build and run commands
└── README.md         # This file
```

## 🛠️ Development

### Available Make Commands

```bash
make build    # Build the project
make run      # Build and run the project
make go-run   # Run without building
make test     # Run tests
make clean    # Clean build artifacts
```

### Running Tests

```bash
make test
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is open source and available under the MIT License.

## 🎯 Roadmap

- [ ] Add support for more Redis-like commands (DEL, EXISTS, etc.)
- [ ] Implement data persistence
- [ ] Add support for different data types (Lists, Sets, Hashes)
- [ ] Implement authentication
- [ ] Add Redis protocol (RESP) support
- [ ] Improve error handling and logging
- [ ] Add comprehensive test coverage

## 👏 Acknowledgments

Inspired by Redis and built as a learning project to explore Go's networking capabilities and concurrent programming patterns.

---

**Note:** This is an educational project and not intended for production use.
