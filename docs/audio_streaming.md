
Improving user experience through audio streaming

# Status

06/19/24: Drafted and Decided

# Context

The implementation prior to this proposed change has a significant input lag
between the user clicking the play button and the delivery of audio. The cause
of this lag is derived from the method of storage and retreival of audio files.
When the user clicks "play", a backend GET call is made to retrieve the audio
file. The files are a 10-100 Mbs in size, and can be larger based on the
recording's duration.

## References
- https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio
- https://developer.mozilla.org/en-US/docs/Web/API/Media_Source_Extensions_API
- https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Range 
- https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/HEAD

# Decision

The audio to be played will be broken up into smaller tracks to be streamed to the front end.

## Frontend Enhancements

The Web UI will utilize the Media Source API. When an audio element is selected
for playback, the corresponding media source object will begin its "sourceopen"
event where it will gradually fetch the audio data from the backend and append
it to the media source object's buffer.

## Backend Enhancements

HEAD http requests will be supported for audio endpoints. This will enable
clients to fetch metadata about the audio file without the need to retrieve its
content. The following metadata will be provided to the frontend:
- chunkSize : The ideal size of chunks for streaming, in bytes. For now this
  value will be arbitrarily set as 1 Mb (This is what ChatGPT recommended.)
- fileSize: This is the total size of the file, in bytes.

The Range http parameter will be supported for GET calls to audio endpoints.
This will enable clients to retreive small chunks of the audio file at a time.
Per the specification, range requests will return 206 on success, and 416 for
invalid range values.

# Consequences

By adopting the described method for audio retreival, users of the audio-server
will be able to play audio files without having to wait for the entire file to
load. This is a substantial improvement to the overall user experience.
