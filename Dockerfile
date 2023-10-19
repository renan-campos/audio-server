# Use the official Go image as the base image
FROM golang:1.20.7
WORKDIR /app
COPY . . 
RUN make
EXPOSE 1323
EXPOSE 1324
CMD ["bin/audio-server"]
