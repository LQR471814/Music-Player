## music-player

> A self-hosted, web-based, nerd-friendly, music player.

### setting up

- `npm install`
- [protoc](https://grpc.io/docs/protoc-installation/) `sudo apt install protobuf-compiler`

to generate protobuf files, run `make protobuf` in the root directory.

to run the server (both golang and vite) you must have a certificate and key in the form of `server/host-crt.pem` and `server/host-key.prm`

to generate such files, you should look into self signing with `openssl` or `xca`

### roadmap

- [x] dynamic themes on the client
- [ ] remove player jankiness
- [ ] adaptive layouts
    - [x] desktop layouts
    - [ ] phone layouts
    - [ ] everything in-between layouts
- [ ] editing
    - [ ] album names
    - [ ] track names
    - [ ] album metadata
    - [ ] track metadata
- [ ] persistent audio caching
- [ ] import audio sources
    - [ ] upload audio
    - [ ] from youtube
        - [ ] without streaming
        - [ ] with streaming
    - [ ] from soundcloud
        - [ ] without streaming
        - [ ] with streaming
- [ ] audio processing
    - [ ] compressor
    - [ ] equalizer
    - [ ] stereo shaper
    - [ ] visualization EQ

### backgrounds

Background 1 by [Anthony Reungère](https://unsplash.com/@anthonyreungere?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText) on [Unsplash](https://unsplash.com/s/photos/city?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText)

Background 2 by [Max Bender](https://unsplash.com/@maxwbender?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText) on [Unsplash](https://unsplash.com/s/photos/city?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText)

