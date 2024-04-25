# Audio Forward

Forward audio from one PC to another using a TCP socket

## Usage

Start the server on the remote PC:

```bash
./audiof server
```

Connect to server:

```bash
./audiof client 192.168.1.105:8322
```

## Todo

  - [ ] Compress data before sending;
  - [ ] Fix some underrun on client;
  - [ ] User can set listening port;
  - [ ] Accept different sample rates.
