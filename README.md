# BBC-English Wechat Mini Program

## About

I used to use this mini-program to provide a convenient way for me to access BBC Learning English. Because in some palaces, we can't visit it directly and need a proxy.

This program contains two parts:

1. The crawler of BBC learning English and uploading audios/videos to Alibaba Cloud OSS
2. providing APIs to get contents from OSS

Thus, to run this program, you may need two servers, one can access BBC and run the crawler; one is used to provide APIs for the mini program. Details are listed below.

## Working Flow

<img src="http://qiniu.lastmayday.com/BBC_English_Mini_Program.png" width="600px">

First, we set a cron to run the crawler on server 1; server 1 can access the BBC directly.
Then, download the audios/videos from BBC, upload them to OSS and get a new OSS link to access these media.
Meanwhile, write the metadata of audios/videos to the database. The metadata contains id, title, transcripts, OSS link, etc.

From server 2, it reads metadata from the database and provides APIs for the mini program.

## Directories functions

- `cmd`: the main function
- `config`: configs of BBC, database, OSS, etc.
- `internal`: the crawler
- `scripts`: scripts to run crawler or APP
- `web`: APIs for APP
- `wx`: Wechat mini program

## Mini program screenshoots

<img src="http://qiniu.lastmayday.com/BBC_English_Screenshot_list.png" width="500px">

<img src="http://qiniu.lastmayday.com/BBC_English_Screenshot_detail.png" width="500px">
