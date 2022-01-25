# YMM

## Features

1. Download music from Youtube
   * Download a video file as audio
   * Download a playlist as audio files
2. Manage downloaded files
   * Make structured directory by creator
   * Put files in its creators' directory
   * Let user change creator directory's name

## Process

### Download a video file as audio

1. User execute command
2. Fetch meta data from youtube
3. Get creator data from DB or create
4. Put downloaded file into creator's directory
5. Execute loudness scanner on downloaded file

### Download a playlist as audio files

1. User execute command
2. Fetch meta data from youtube
3. Execute [Download a video file as audio](#Download-a-video-file-audio) for each video

## Data

### Creator

| Field      | Type   | Example    |
| ---------- | ------ | ---------- |
| creator_id | String | MARETU     |
| directory  | String | /what/ever |

### Audio

| Field      | Type    | Example  |
| :--------- | :------ | -------- |
| creator_id | String  | MARETU   |
| filename   | String  | asdf.mp3 |
| overrided  | Boolean | true     |

