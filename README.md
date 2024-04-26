# Audio Forward

Forward audio from one PC to another using a TCP socket

## Build

Dependencies:
  - portaudio19-dev

Compile the program:

```bash
go build .
```

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
  - [X] Fix some underrun on client(tested with a video from Akita on Rails and some musics);
  - [ ] User can set listening port;
  - [ ] Accept different sample rates.
