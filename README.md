# overlayer

Try to run this app:

```bash
podman build -t overlayer-dev . && \
podman run \
  -v $(pwd)/logo.svg:/app/logo.svg:ro,z \
  -e OVERLAY_LOGO_PATH="logo.svg" \
  -e OVERLAY_LOGO_HEIGHT="128" \
  -e OVERLAY_LOGO_WIDTH="128" \
  -e OVERLAY_LOGO_ALPHA="1" \
  -e STREAMS_0_SRC="rtmp://example.com/stream" \
  -e STREAMS_0_DST="srt://example.com/live" \
  overlayer-dev
```

Example of ffmpeg cli:

```bash
ffmpeg -i "INPUT" -i logo.svg -filter_complex "[1:v]scale=128:128,format=rgba,colorchannelmixer=aa=0.75[logo];[0:v][logo]overlay=24:(main_h-overlay_h)/2" -c:v libx264 -preset ultrafast -c:a aac -f flv "OUTPUT"
```

Example of ffmpeg Dockerfile:

```bash
podman run --rm -it \
  -v $(pwd)/logo.svg:/logo.svg:ro,z \
  localhost/ffmpeg \
  -re \
  -i "INPUT" \
  -i "/logo.svg" \
  -filter_complex "[1:v]scale=128:128,format=rgba,colorchannelmixer=aa=0.75[logo];[0:v][logo]overlay=24:(main_h-overlay_h)/2" \
  -c:v "libx264" \
  -preset "ultrafast" \
  -tune "zerolatency" \
  -c:a "aac" \
  -flush_packets 0 \
  -max_muxing_queue_size 2048 \
  -rtbufsize 1500k \
  -fifo_size 100000 \
  -reconnect 1 \
  -reconnect_streamed 1 \
  -reconnect_delay_max 5 \
  -analyzeduration 50000000 \
  -probesize 100000000 \
  -err_detect "ignore_err" \
  -loglevel "warning" \
  -f "flv" "OUTPUT"
```
